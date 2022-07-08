# Measurement tool for firehose performance

## Build

```shell
go build 
```

## Run

```shell
export STREAMINGFAST_API_KEY=...
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
  -head
        Ignores start-block and block-range settings and tests live blocks only
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
./sf-perf -connections 3 -hosts eth.firehose.pinax.network:9000 -start-block 1000000 -block-range 10000    
2022-07-08T14:03:36.306+0200 INFO (sf-perf) initialised worker {"id": 0, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1000000  stop_block_num:1010000  fork_steps:STEP_NEW"}
2022-07-08T14:03:36.641+0200 INFO (sf-perf) initialised worker {"id": 1, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1010000  stop_block_num:1020000  fork_steps:STEP_NEW"}
2022-07-08T14:03:36.966+0200 INFO (sf-perf) initialised worker {"id": 2, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1020000  stop_block_num:1030000  fork_steps:STEP_NEW"}
2022-07-08T14:03:36.966+0200 INFO (sf-perf) starting measurement...
2022-07-08T14:03:58.509+0200 INFO (sf-perf) done
2022-07-08T14:03:58.509+0200 INFO (sf-perf) worker results {"worker_id": 0, "start_block": 1000000, "stop_block": 1010000, "num_blocks": 10001, "total_time": "13.859068499s", "time_to_first_block": "326.131µs", "total_size": "61 MB", "blocks_per_second": 721, "bytes_per_second": "4.4 MB"}
2022-07-08T14:03:58.509+0200 INFO (sf-perf) worker results {"worker_id": 1, "start_block": 1010000, "stop_block": 1020000, "num_blocks": 10001, "total_time": "21.541808788s", "time_to_first_block": "262.917µs", "total_size": "76 MB", "blocks_per_second": 464, "bytes_per_second": "3.5 MB"}
2022-07-08T14:03:58.509+0200 INFO (sf-perf) worker results {"worker_id": 2, "start_block": 1020000, "stop_block": 1030000, "num_blocks": 10001, "total_time": "15.115675526s", "time_to_first_block": "201.370546ms", "total_size": "59 MB", "blocks_per_second": 661, "bytes_per_second": "3.9 MB"}
2022-07-08T14:03:58.509+0200 INFO (sf-perf) result summary {"num_workers": 3, "num_blocks": 30003, "avg_time": "16.838850937s", "avg_time_to_first_block": "67.319864ms", "total_size": "196 MB", "total_blocks_per_second": 1847, "avg_blocks_per_second": 615, "total_bytes_per_second": "12 MB", "avg_bytes_per_second": "3.9 MB"}
```