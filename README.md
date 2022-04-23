This is a simple dummy tool for Google Pub/Sub that can monitor a set of topics (subscribe to them) and save the messages they receive for retrospective analysis.
Among others, this can be useful for security research purposes/blackbox testing to look under the hood of a service that leverages Pub/Sub for control.

By default, the tool queries the topics present in the specified project and subscribes to them all. You can fine tune this behaviour by specifying
the list of topic IDs you would like to monitor. Each message received will be persisted to the directory specified via the -log-dir command line parameter.

For example, when the Cloud Storage Transfer service is carrying out a copy operation, the output directory is populated with files like this:

```
-rw-r--r-- 1 root root   93 Apr 23 21:45 0061-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  148 Apr 23 21:45 0062-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:45 0063-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  148 Apr 23 21:46 0064-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:46 0065-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  148 Apr 23 21:46 0066-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:46 0067-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  384 Apr 23 21:46 0068-cloud-ingest-start.bin
-rw-r--r-- 1 root root  251 Apr 23 21:46 0069-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:46 0070-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  251 Apr 23 21:46 0071-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:46 0072-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  251 Apr 23 21:46 0073-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:46 0074-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  251 Apr 23 21:46 0075-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:46 0076-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  251 Apr 23 21:47 0077-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   93 Apr 23 21:47 0078-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  251 Apr 23 21:47 0079-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root  597 Apr 23 21:47 0080-p-storage-transfer-1-cloud-ingest-list.bin
-rw-r--r-- 1 root root   93 Apr 23 21:47 0081-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root 1393 Apr 23 21:47 0082-cloud-ingest-list-progress.bin
-rw-r--r-- 1 root root  251 Apr 23 21:47 0083-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root  716 Apr 23 21:47 0084-cloud-ingest-process-list.bin
-rw-r--r-- 1 root root  451 Apr 23 21:47 0085-p-storage-transfer-1-cloud-ingest-copy.bin
-rw-r--r-- 1 root root   96 Apr 23 21:47 0086-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  843 Apr 23 21:47 0087-cloud-ingest-copy-progress.bin
-rw-r--r-- 1 root root  148 Apr 23 21:47 0088-p-storage-transfer-1-cloud-ingest-control.bin
-rw-r--r-- 1 root root   99 Apr 23 21:47 0089-cloud-ingest-pulse.bin
-rw-r--r-- 1 root root  148 Apr 23 21:47 0090-p-storage-transfer-1-cloud-ingest-control.bin
```

Each of these files contain the raw message data blob, which is a protocol buffer in this case:

```
protoc --decode_raw <0063-cloud-ingest-pulse.bin
1 {
  1: "5ad9c1e1be6f"
  2: "1338"
}
3: "3.0.2392"
4: "/tmp"
7: 1290000
20 {
  2: 2
}
21: 1
22 {
  1: 3
}
22 {
  1: 1
}
22 {
  1: 2
}
22 {
  1: 4
}
24: "p-storage-transfer-1"
```
