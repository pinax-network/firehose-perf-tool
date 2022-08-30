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
        Authentication endpoint to retrieve access tokens from. (default "https://auth.eosnation.io")
  -block-range int
        Block range for each connection (default 10000)
  -compression
        apply gzip compression on the grpc connection
  -connections int
        Number of concurrent connections to measure (default 10)
  -eth-call-filter-multi string
        Advanced filter. List of 'address[+address[+...]]:eventsig[+eventsig[+...]]' pairs, ex: 'dead+beef:1234+5678,:0x44,0x12:' results in 3 filters.
  -eth-log-filter-multi string
        Advanced filter. List of 'address[+address[+...]]:eventsig[+eventsig[+...]]' pairs, ex: 'dead+beef:1234+5678,:0x44,0x12:' results in 3 filters.
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
./sf-perf -connections 3 -hosts eth.firehose.pinax.network:9000 -block-range 100 -eth-log-filter-multi 0x3a8778A58993bA4B941f85684D74750043A4bB5f:0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
2022-08-30T15:53:53.457+0200 INFO (sf-perf) initialised worker {"id": 0, "host": "eth.firehose.pinax.network:9000", "request_options": "stop_block_num:100 fork_steps:STEP_NEW transforms:{[type.googleapis.com/sf.ethereum.transform.v1.MultiLogFilter]:{log_filters:{addresses:\":\\x87x\\xa5\\x89\\x93\\xbaK\\x94\\x1f\\x85hMtu\\x00C\\xa4\\xbb_\" event_signatures:\"\\xdd\\xf2R\\xad\\x1b\\xe2ți°h\\xfc7\\x8d\\xaa\\x95+\\xa7\\xf1cġ\\x16(\\xf5ZM\\xf5#\\xb3\\xef\"}}}"}
2022-08-30T15:53:54.061+0200 INFO (sf-perf) initialised worker {"id": 1, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:100 stop_block_num:200 fork_steps:STEP_NEW transforms:{[type.googleapis.com/sf.ethereum.transform.v1.MultiLogFilter]:{log_filters:{addresses:\":\\x87x\\xa5\\x89\\x93\\xbaK\\x94\\x1f\\x85hMtu\\x00C\\xa4\\xbb_\" event_signatures:\"\\xdd\\xf2R\\xad\\x1b\\xe2ți°h\\xfc7\\x8d\\xaa\\x95+\\xa7\\xf1cġ\\x16(\\xf5ZM\\xf5#\\xb3\\xef\"}}}"}
2022-08-30T15:53:54.408+0200 INFO (sf-perf) initialised worker {"id": 2, "host": "eth.firehose.pinax.network:9000", "request_options": "start_block_num:200 stop_block_num:300 fork_steps:STEP_NEW transforms:{[type.googleapis.com/sf.ethereum.transform.v1.MultiLogFilter]:{log_filters:{addresses:\":\\x87x\\xa5\\x89\\x93\\xbaK\\x94\\x1f\\x85hMtu\\x00C\\xa4\\xbb_\" event_signatures:\"\\xdd\\xf2R\\xad\\x1b\\xe2ți°h\\xfc7\\x8d\\xaa\\x95+\\xa7\\xf1cġ\\x16(\\xf5ZM\\xf5#\\xb3\\xef\"}}}"}
2022-08-30T15:53:54.408+0200 INFO (sf-perf) starting measurement...
2022-08-30T15:54:17.799+0200 INFO (sf-perf) done
2022-08-30T15:54:17.800+0200 INFO (sf-perf) worker results {"worker_id": 0, "start_block": 0, "stop_block": 100, "num_blocks": 100, "total_time": "15.06216797s", "time_to_first_block": "5.905638024s", "total_size": "107 kB", "blocks_per_second": 6, "bytes_per_second": "7.1 kB"}
2022-08-30T15:54:17.800+0200 INFO (sf-perf) worker results {"worker_id": 1, "start_block": 100, "stop_block": 200, "num_blocks": 101, "total_time": "21.725414517s", "time_to_first_block": "14.866380202s", "total_size": "112 kB", "blocks_per_second": 4, "bytes_per_second": "5.1 kB"}
2022-08-30T15:54:17.800+0200 INFO (sf-perf) worker results {"worker_id": 2, "start_block": 200, "stop_block": 300, "num_blocks": 101, "total_time": "23.391072223s", "time_to_first_block": "15.651838457s", "total_size": "115 kB", "blocks_per_second": 4, "bytes_per_second": "4.9 kB"}
2022-08-30T15:54:17.800+0200 INFO (sf-perf) result summary {"num_workers": 3, "failed_workers": 0, "num_blocks": 302, "avg_time": "20.05955157s", "avg_time_to_first_block": "12.141285561s", "total_size": "334 kB", "total_blocks_per_second": 15, "avg_blocks_per_second": 5, "total_bytes_per_second": "17 kB", "avg_bytes_per_second": "5.7 kB"}
2022-08-30T15:54:17.800+0200 INFO (sf-perf) input parameters {"start-block": 0, "block-range": 100, "connections": 3, "insecure": false, "plaintext": false, "hosts": "eth.firehose.pinax.network:9000", "head": false, "auth-endpoint": "https://auth.pinax.network", "eth-log-filter-multi": "0x3a8778A58993bA4B941f85684D74750043A4bB5f:0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "eth-call-filter-multi": ""}
```