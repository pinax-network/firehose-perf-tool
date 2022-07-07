package measurement

import (
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
)

func PrintResults(measurements []*Measurement) {

	for _, m := range measurements {

		totalTime := m.Blocks[len(m.Blocks)-1].BlockReceivedAt.Sub(m.StartTime)
		timeToFirstBlock := m.Blocks[0].BlockReceivedAt.Sub(m.StartTime)
		// streamTime := m.Blocks[len(m.Blocks)-1].BlockReceivedAt.Sub(m.Blocks[0].BlockReceivedAt)

		numBlocks := len(m.Blocks)
		totalBlockSize := 0

		for _, b := range m.Blocks {
			totalBlockSize += b.BlockSize
		}

		blockThroughput := float64(len(m.Blocks)) / totalTime.Seconds()
		sizeThroughput := float64(totalBlockSize) / totalTime.Seconds()

		zlog.Info("worker results",
			zap.Int("worker_id", m.WorkerId),
			zap.Int64("start_block", m.RequestOptions.StartBlockNum),
			zap.Uint64("stop_block", m.RequestOptions.StopBlockNum),
			zap.Int("num_blocks", numBlocks),
			zap.Duration("total_time", totalTime),
			zap.Duration("time_to_first_block", timeToFirstBlock),
			// zap.Duration("stream_time", streamTime),
			zap.String("total_size", humanize.Bytes(uint64(totalBlockSize))),
			zap.Int("blocks_per_second", int(blockThroughput)),
			zap.String("bytes_per_second", humanize.Bytes(uint64(sizeThroughput))),
		)
	}
}
