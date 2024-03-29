package measurement

import (
	"errors"
	pbfirehose "github.com/streamingfast/pbgo/sf/firehose/v1"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
	"time"
)

type Worker struct {
	id          int
	wg          *sync.WaitGroup
	stream      pbfirehose.Stream_BlocksClient
	measurement *Measurement
	shutdown    *atomic.Bool
}

type Measurement struct {
	WorkerId       int
	HasFailed      bool
	RequestOptions *pbfirehose.Request
	StartTime      time.Time
	Blocks         []BlockResult
}

type BlockResult struct {
	EstimatedBlockNum int
	BlockSize         int
	BlockReceivedAt   time.Time
}

func NewWorker(id int, wg *sync.WaitGroup, stream pbfirehose.Stream_BlocksClient, requestOptions *pbfirehose.Request) *Worker {
	return &Worker{
		id:       id,
		wg:       wg,
		stream:   stream,
		shutdown: atomic.NewBool(false),
		measurement: &Measurement{
			WorkerId:       id,
			HasFailed:      false,
			RequestOptions: requestOptions,
			Blocks:         make([]BlockResult, 0),
		},
	}
}

func (n *Worker) StartMeasurement() {

	n.measurement.StartTime = time.Now()
	estimatedBlockNum := int(n.measurement.RequestOptions.StartBlockNum)

	for {
		response, err := n.stream.Recv()
		if errors.Is(err, io.EOF) || n.shutdown.Load() {
			// we are done here
			n.wg.Done()
			return
		} else if err != nil {
			n.measurement.HasFailed = true
			zlog.Error("measurement failed", zap.Int("worker_id", n.id), zap.Error(err))
			n.wg.Done()
			return
		}

		blockResult := BlockResult{
			EstimatedBlockNum: estimatedBlockNum,
			BlockSize:         proto.Size(response),
			BlockReceivedAt:   time.Now(),
		}

		n.measurement.Blocks = append(n.measurement.Blocks, blockResult)
		// we don't know the exact block number as we don't parse the result, but we estimate by incrementing from the start block
		estimatedBlockNum++
	}
}

func (n *Worker) StopMeasurement() {
	n.shutdown.Store(true)
}

func (n *Worker) GetResults() *Measurement {
	return n.measurement
}
