package measurement

import (
	pbfirehose "github.com/streamingfast/pbgo/sf/firehose/v1"
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
}

type Measurement struct {
	WorkerId       int
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
		id:     id,
		wg:     wg,
		stream: stream,
		measurement: &Measurement{
			WorkerId:       id,
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
		if err == io.EOF {
			// we are done here
			n.wg.Done()
			return
		} else if err != nil {
			zlog.Error("measurement failed", zap.Int("worker_id", n.id), zap.Error(err))
			return
		}

		blockResult := BlockResult{
			EstimatedBlockNum: estimatedBlockNum,
			BlockSize:         proto.Size(response),
			BlockReceivedAt:   time.Now(),
		}

		n.measurement.Blocks = append(n.measurement.Blocks, blockResult)
	}
}

func (n *Worker) GetResults() *Measurement {
	return n.measurement
}
