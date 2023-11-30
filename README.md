# Key-Value Store with OrderedMap Integration

## Overview

Welcome to the Key-Value Store project! This project aims to provide a persistent key-value storage solution with a simple and efficient HTTP API. The architecture is inspired by LSM tree models, and the implementation draws inspiration from well-established projects like leveldb and rocksdb.

## Features

- **HTTP API:** Exposes three crucial endpoints for key-value operations:
  - `GET http://localhost:8080/get?key=keyName`: Retrieve the value associated with the provided key.
  - `POST http://localhost:8080/set`: Set a key-value pair. Send the data in JSON format in the request body.
  - `DELETE http://localhost:8080/del?key=keyName`: Delete a key from the store and return its existing value if present.

- **Write-Ahead Log (WAL):** Ensures crash safety by writing all write operations first to a memtable (sorted map of key-value pairs) and appending them to the WAL before responding to the user.

- **SST Files (Sorted String Table):** Implements a mechanism for periodically flushing the contents of the memtable to disk as an SST file. This process helps in maintaining a snapshot of the memtable on disk and preventing the number of SST files from becoming too large.

- **OrderedMap Integration:** Utilizes the [OrderedMap](https://github.com/example/orderedmap) library, a sorted map data structure, to enhance the efficiency of key-value storage and retrieval.

## Project Structure

- **MemDB:**
  - Handles in-memory storage of key-value pairs.
  - Manages the memtable, write-ahead log (WAL), and the flushing mechanism.

- **HTTP API:**
  - Listens for incoming HTTP requests and delegates them to the MemDB for processing.
  - Provides a user-friendly interface for interacting with the key-value store.

- **SST Files:**
  - Crucial for persistence and durability.
  - Contains functions for parsing existing SST files and flushing memtable contents to create new SST files.

- **OrderedMap:**
  - Integrated for optimized key ordering and retrieval.

## Getting Started

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/mohiZzine/key-value-store.git
  ```