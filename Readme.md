# Measurement tool for firehose performance

## Run

```shell
export STREAMINGFAST_API_KEY=...
go build 
./sf-perf 
```

### Parameters 

```shell
Usage of ./sf-perf:
  -auth-endpoint string
         (default "https://auth.eosnation.io")
  -block-range int
        Block range for each connection (default 10000)
  -connections int
        Number of concurrent connections to measure (default 10)
  -hosts string
        Comma separated list of hosts
  -insecure
        Skip TLS certificate verification
  -plaintext
        Use plaintext connection
  -start-block int
        Start block to start the measuring from
```

### Output

```shell
./sf-perf -connections 3 -hosts eth.firehose.pinax.network:9000 -start-block 1000000 -block-range 100000
2022-07-07T15:29:32.627+0200 INFO (sf) initialised worker {"id": 0, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1000000 stop_block_num:1100000 fork_steps:STEP_NEW"}
2022-07-07T15:29:33.010+0200 INFO (sf) initialised worker {"id": 1, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1100000 stop_block_num:1200000 fork_steps:STEP_NEW"}
2022-07-07T15:29:33.373+0200 INFO (sf) initialised worker {"id": 2, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1200000 stop_block_num:1300000 fork_steps:STEP_NEW"}
2022-07-07T15:29:33.373+0200 INFO (sf) starting measurement...
2022-07-07T15:33:02.778+0200 INFO (sf) done
2022-07-07T15:33:02.778+0200 INFO (sf) worker results {"worker_id": 0, "start_block": 1000000, "stop_block": 1100000, "num_blocks": 100001, "total_time": "2m17.396076258s", "time_to_first_block": "274.362µs", "stream_time": "2m17.395801896s", "total_size": "621 MB", "blocks_per_second": 727, "bytes_per_second": "4.5 MB"}
2022-07-07T15:33:02.779+0200 INFO (sf) worker results {"worker_id": 1, "start_block": 1100000, "stop_block": 1200000, "num_blocks": 100001, "total_time": "2m54.383885664s", "time_to_first_block": "186.895µs", "stream_time": "2m54.383698769s", "total_size": "764 MB", "blocks_per_second": 573, "bytes_per_second": "4.4 MB"}
2022-07-07T15:33:02.779+0200 INFO (sf) worker results {"worker_id": 2, "start_block": 1200000, "stop_block": 1300000, "num_blocks": 100001, "total_time": "3m29.403862025s", "time_to_first_block": "231.28854ms", "stream_time": "3m29.172573485s", "total_size": "1.0 GB", "blocks_per_second": 478, "bytes_per_second": "5.0 MB"}
```