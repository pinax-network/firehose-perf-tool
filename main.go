package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	dfuse "github.com/streamingfast/client-go"
	"github.com/streamingfast/dgrpc"
	"github.com/streamingfast/logging"
	pbfirehose "github.com/streamingfast/pbgo/sf/firehose/v1"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"os"
	"sf-perf/measurement"
	"strings"
	"sync"
)

func init() {
	logging.ApplicationLogger("sf", "github.com/pinax-network/firehose-perf-tool")
}

func main() {

	startBlockPtr := flag.Int("start-block", 0, "")
	blockRangePtr := flag.Int("block-range", 10000, "")
	connectionsPtr := flag.Int("connections", 10, "")
	hostsPtr := flag.String("hosts", "", "Comma separated list of hosts")
	authEndpoint := flag.String("auth-endpoint", "https://auth.eosnation.io", "")

	flag.Parse()

	workerPool := make([]*measurement.Worker, *connectionsPtr)
	hosts := strings.Split(*hostsPtr, ",")
	wg := &sync.WaitGroup{}

	zlog.Info("connections pointer", zap.Any("connectionsPtr", connectionsPtr))

	// init workers
	for i := 0; i < *connectionsPtr; i++ {

		workerEndpoint := hosts[i%len(hosts)]
		requestOptions := &pbfirehose.Request{
			StartBlockNum: int64(*startBlockPtr + (i * *blockRangePtr)),
			StopBlockNum:  uint64(*startBlockPtr + (i * *blockRangePtr) + *blockRangePtr),
			ForkSteps:     []pbfirehose.ForkStep{pbfirehose.ForkStep_STEP_NEW},
		}
		workerStream, err := newStream(context.Background(), *authEndpoint, workerEndpoint, requestOptions)
		if err != nil {
			zlog.Fatal("failed to initialise stream", zap.Error(err))
		}
		wg.Add(1)

		worker := measurement.NewWorker(i, wg, workerStream, requestOptions)
		workerPool[i] = worker

		zlog.Info("initialised worker", zap.Int("id", i), zap.String("host", workerEndpoint), zap.Any("request_options", requestOptions))
	}

	for _, w := range workerPool {
		go w.StartMeasurement()
	}

	wg.Wait()

	zlog.Info("finished measurement")
	measurements := make([]*measurement.Measurement, len(workerPool))
	for i, w := range workerPool {
		measurements[i] = w.GetResults()
	}

	measurement.PrintResults(measurements)
}

func newStream(ctx context.Context, authEndpoint, endpoint string, requestOptions *pbfirehose.Request) (stream pbfirehose.Stream_BlocksClient, err error) {

	var clientOptions []dfuse.ClientOption
	skipAuth := false
	apiKey := os.Getenv("STREAMINGFAST_API_KEY")
	if apiKey == "" {
		clientOptions = []dfuse.ClientOption{dfuse.WithoutAuthentication()}
		skipAuth = true
	}

	if viper.GetBool("skip-auth") {
		clientOptions = []dfuse.ClientOption{dfuse.WithoutAuthentication()}
		skipAuth = true
	}

	if authEndpoint != "" && !skipAuth {
		clientOptions = []dfuse.ClientOption{dfuse.WithAuthURL(authEndpoint)}
	}

	client, err := dfuse.NewClient(endpoint, apiKey, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to create streamingfast client")
	}

	useInsecureTSLConnection := viper.GetBool("global-insecure")
	usePlainTextConnection := viper.GetBool("global-plaintext")

	if useInsecureTSLConnection && usePlainTextConnection {
		return nil, fmt.Errorf("option --insecure and --plaintext are mutually exclusive, they cannot be both specified at the same time")
	}

	var dialOptions []grpc.DialOption
	switch {
	case usePlainTextConnection:
		zlog.Debug("Configuring transport to use a plain text connection")
		dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	case useInsecureTSLConnection:
		zlog.Debug("Configuring transport to use an insecure TLS connection (skips certificate verification)")
		dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))}
	}

	conn, err := dgrpc.NewExternalClient(endpoint, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to create external gRPC client")
	}

	grpcCallOpts := []grpc.CallOption{}

	if !skipAuth {
		tokenInfo, err := client.GetAPITokenInfo(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve StreamingFast API token: %w", err)
		}
		rpcCredentials := oauth.NewOauthAccess(&oauth2.Token{AccessToken: tokenInfo.Token, TokenType: "Bearer"})
		grpcCallOpts = append(grpcCallOpts, grpc.PerRPCCredentials(rpcCredentials))
	}

	firehoseClient := pbfirehose.NewStreamClient(conn)

	zlog.Debug("Initiating stream with remote endpoint", zap.String("endpoint", endpoint))
	stream, err = firehoseClient.Blocks(context.Background(), requestOptions, grpcCallOpts...)
	if err != nil {
		return nil, fmt.Errorf("unable to start blocks stream: %w", err)
	}

	return
}
