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
	"os/signal"
	"sf-perf/measurement"
	"strings"
	"sync"
	"syscall"
)

func init() {
	logging.ApplicationLogger("sf-perf", "github.com/pinax-network/firehose-perf-tool")
}

func main() {

	startBlockPtr := flag.Int("start-block", 0, "Start block to start the measuring from")
	blockRangePtr := flag.Int("block-range", 10000, "Block range for each connection")
	connectionsPtr := flag.Int("connections", 10, "Number of concurrent connections to measure")
	insecurePtr := flag.Bool("insecure", false, "Skip TLS certificate verification")
	plaintextPtr := flag.Bool("plaintext", false, "Use plaintext connection")
	hostsPtr := flag.String("hosts", "", "Comma separated list of hosts")
	headPtr := flag.Bool("head", false, "Ignores start-block and block-range settings and tests live blocks only")
	authEndpoint := flag.String("auth-endpoint", "https://auth.eosnation.io", "")

	flag.Parse()

	workerPool := make([]*measurement.Worker, *connectionsPtr)
	hosts := strings.Split(*hostsPtr, ",")
	wg := &sync.WaitGroup{}

	// init workers
	for i := 0; i < *connectionsPtr; i++ {

		startBlockNum := int64(*startBlockPtr + (i * *blockRangePtr))
		stopBlockNum := uint64(*startBlockPtr + (i * *blockRangePtr) + *blockRangePtr)

		if *headPtr {
			startBlockNum = -1
			stopBlockNum = 0xFFFFFFFFFFFFFFFF
		}

		workerEndpoint := hosts[i%len(hosts)]
		requestOptions := &pbfirehose.Request{
			StartBlockNum: startBlockNum,
			StopBlockNum:  stopBlockNum,
			ForkSteps:     []pbfirehose.ForkStep{pbfirehose.ForkStep_STEP_NEW},
		}
		workerStream, err := newStream(context.Background(), *authEndpoint, workerEndpoint, *insecurePtr, *plaintextPtr, requestOptions)
		if err != nil {
			zlog.Fatal("failed to initialise stream", zap.Error(err))
		}
		wg.Add(1)

		worker := measurement.NewWorker(i, wg, workerStream, requestOptions)
		workerPool[i] = worker

		zlog.Info("initialised worker", zap.Int("id", i), zap.String("host", workerEndpoint), zap.Any("request_options", requestOptions))
	}

	zlog.Info("starting measurement...")

	for _, w := range workerPool {
		go w.StartMeasurement()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		zlog.Info("notifying workers about shutdown...")
		for _, w := range workerPool {
			w.StopMeasurement()
		}
	}()

	wg.Wait()

	zlog.Info("done")
	measurements := make([]*measurement.Measurement, len(workerPool))
	for i, w := range workerPool {
		measurements[i] = w.GetResults()
	}

	measurement.PrintResults(measurements)
}

func newStream(ctx context.Context, authEndpoint, endpoint string, insecureConn, plaintextConn bool, requestOptions *pbfirehose.Request) (stream pbfirehose.Stream_BlocksClient, err error) {

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

	if insecureConn && plaintextConn {
		return nil, fmt.Errorf("option --insecure and --plaintext are mutually exclusive, they cannot be both specified at the same time")
	}

	var dialOptions []grpc.DialOption
	switch {
	case plaintextConn:
		zlog.Debug("Configuring transport to use a plain text connection")
		dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	case insecureConn:
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
