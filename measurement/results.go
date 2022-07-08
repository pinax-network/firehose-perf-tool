package measurement

import (
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
	"time"
)

func PrintResults(measurements []*Measurement) {

	overallNumBlocks := 0
	overallBlockSize := 0
	overallTime := time.Duration(0)
	overallTimeToFirstBlock := time.Duration(0)
	overallBlockThroughput := float64(0)
	overallSizeThroughput := float64(0)

	for _, m := range measurements {

		totalTime := m.Blocks[len(m.Blocks)-1].BlockReceivedAt.Sub(m.StartTime)
		overallTime += totalTime

		timeToFirstBlock := m.Blocks[0].BlockReceivedAt.Sub(m.StartTime)
		overallTimeToFirstBlock += timeToFirstBlock

		numBlocks := len(m.Blocks)
		totalBlockSize := 0

		for _, b := range m.Blocks {
			totalBlockSize += b.BlockSize
		}

		overallNumBlocks += numBlocks
		overallBlockSize += totalBlockSize

		blockThroughput := float64(len(m.Blocks)) / totalTime.Seconds()
		overallBlockThroughput += blockThroughput

		sizeThroughput := float64(totalBlockSize) / totalTime.Seconds()
		overallSizeThroughput += sizeThroughput

		zlog.Info("worker results",
			zap.Int("worker_id", m.WorkerId),
			zap.Int64("start_block", m.RequestOptions.StartBlockNum),
			zap.Uint64("stop_block", m.RequestOptions.StopBlockNum),
			zap.Int("num_blocks", numBlocks),
			zap.Duration("total_time", totalTime),
			zap.Duration("time_to_first_block", timeToFirstBlock),
			zap.String("total_size", humanize.Bytes(uint64(totalBlockSize))),
			zap.Int("blocks_per_second", int(blockThroughput)),
			zap.String("bytes_per_second", humanize.Bytes(uint64(sizeThroughput))),
		)
	}

	zlog.Info("result summary",
		zap.Int("num_workers", len(measurements)),
		zap.Int("num_blocks", overallNumBlocks),
		zap.Duration("avg_time", overallTime/time.Duration(len(measurements))),
		zap.Duration("avg_time_to_first_block", overallTimeToFirstBlock/time.Duration(len(measurements))),
		zap.String("total_size", humanize.Bytes(uint64(overallBlockSize))),
		zap.Int("total_blocks_per_second", int(overallBlockThroughput)),
		zap.Int("avg_blocks_per_second", int(overallBlockThroughput/float64(len(measurements)))),
		zap.String("total_bytes_per_second", humanize.Bytes(uint64(overallSizeThroughput))),
		zap.String("avg_bytes_per_second", humanize.Bytes(uint64(overallSizeThroughput/float64(len(measurements))))),
	)
}
