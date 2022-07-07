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
2022-07-07T15:50:44.911+0200 INFO (sf-perf) initialised worker {"id": 0, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1000000  stop_block_num:1100000  fork_steps:STEP_NEW"}
2022-07-07T15:50:45.276+0200 INFO (sf-perf) initialised worker {"id": 1, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1100000  stop_block_num:1200000  fork_steps:STEP_NEW"}
2022-07-07T15:50:45.652+0200 INFO (sf-perf) initialised worker {"id": 2, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:1200000  stop_block_num:1300000  fork_steps:STEP_NEW"}
2022-07-07T15:50:45.652+0200 INFO (sf-perf) starting measurement...
2022-07-07T15:54:15.207+0200 INFO (sf-perf) done
2022-07-07T15:54:15.207+0200 INFO (sf-perf) worker results {"worker_id": 0, "start_block": 1000000, "stop_block": 1100000, "num_blocks": 100001, "total_time": "2m36.855526504s", "time_to_first_block": "174.167µs", "total_size": "621 MB", "blocks_per_second": 637, "bytes_per_second": "4.0 MB"}
2022-07-07T15:54:15.207+0200 INFO (sf-perf) worker results {"worker_id": 1, "start_block": 1100000, "stop_block": 1200000, "num_blocks": 100001, "total_time": "2m42.570311403s", "time_to_first_block": "164.919µs", "total_size": "764 MB", "blocks_per_second": 615, "bytes_per_second": "4.7 MB"}
2022-07-07T15:54:15.208+0200 INFO (sf-perf) worker results {"worker_id": 2, "start_block": 1200000, "stop_block": 1300000, "num_blocks": 100001, "total_time": "3m29.549296382s", "time_to_first_block": "180.970628ms", "total_size": "1.0 GB", "blocks_per_second": 477, "bytes_per_second": "5.0 MB"}
```