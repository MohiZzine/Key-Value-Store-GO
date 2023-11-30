# Key-Value Store with Simple HTTP API

## Introduction

This project implements a persistent key-value store with a simple HTTP API. It exposes the following endpoints:

* GET http://localhost:8081/get?key=keyName: Retrieves the value associated with the specified key.
* POST http://localhost:8081/set: Sets the value associated with the specified key. The key-value pair is provided in the request body as JSON.
* DELETE http://localhost:8081/del?key=keyName: Deletes the specified key and returns its associated value.

The key-value store operates on the LSM tree model for efficient data reading and writing. Write operations are initially stored in the memtable, a sorted map of key-value pairs. Periodically, the memtable is flushed to disk as an SST file (Sorted String Table). To prevent an excessive number of SST files, compaction merges smaller files into larger ones. This compaction process is executed concurrently using a Go routine.

The SST files are in binary format and include the following fields:

* Magic Number: The unique identifier for the application.
* Entry Count: Number of the key-value pairs in the SST File.
* Version: A version number to manage updates to the value.
* Key: The unique identifier for the value.
* Key: The unique identifier for the value.
* Value: The data associated with the key.
* Checksum: A hash value to detect corrupted files.

## Added Dependencies

In this project, the [orderedmap](https://github.com/iancoleman/orderedmap/tree/master) package has been integrated to efficiently manage the ordering of keys in the memtable. This package provides a reliable and performant ordered map implementation.

To integrate this dependency into your project, ensure that Go is installed and execute:


```bash
go get -u github.com/iancoleman/sortedmap
```


## Extras

* Compression: SST files are compressed to save disk space.

## Problem Encountered - Wal Cleaning

Initially, attempts were made to implement log file management with a watermark. Unfortunately, this decision proved to be challenging, resulting in persistent Access Denied Errors, particularly when renaming the temporary file to the log (known to be atomic on Unix-based systems but not on Windows). After extensive debugging efforts over three days, the decision was made to truncate the log file post-flushing. While not the most cost-effective approach, it provides a temporary solution until Wal Cleaning is implemented correctly in future updates.


## Future Improvements

* *Ensuring Atomicity:* Investigate methods to ensure the atomicity of creating SST files during flushing, particularly focusing on the file writing process dependent on the operating system.
* *Concurrent Distributed Database:* Develop a concurrent distributed database to handle multiple clients and achieve high availability.
* *Performance Enhancement:* Explore techniques, such as leveraging Goroutines for parallel processing and optimizing data structures, to enhance the key-value store's performance.

## Getting Started

To run the key-value store, follow these steps:

1. Clone the repository.
2. Start the server: go run .

You can then access the key-value store using aforementioned API endpoints.