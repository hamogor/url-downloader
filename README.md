
# URL Submission and Ranking API

## Overview
This project implements an API to handle URL submissions, track download successes/failures, and provide endpoints to get the top URLs based on the amount of submissions or the latest N submitted. The URLs and their associated data are stored in an in memory map with a doubly linked list for fast accessing for updates as well as maintaining order to query the latest URLs quickly

## Features
- **Submit URL**: Allows submitting a URL and tracks its success/failure.
- **Top URLs**: Fetches the top N URLs based on different sorting criteria like count or latest submitted.
- **Filtering**: Supports sorting and limiting the number of URLs returned.

## Endpoints

### 1. **Submit URL**
- **Endpoint**: `/submit-url`
- **Method**: `POST`
- **Description**: Accepts a URL parameter and records its submission. Tracks the number of successes and failures. If a URL has successfully been fetched in a previous request then the count and success/failure and LastSubmitted time will be updated. If this is the first time that specific URL has been submitted and the GET request to fetch it fails. It will not be stored.
- **Request Body** (JSON):
  ```json
  {
    "url": "http://example.com"
  }
  ```
- **Response**:
  ```json
  {
    "message": "url submitted"
  }
  ```
- **Example**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"url": "http://example.com"}' http://host/submit-url
  ```

### 2. **Get Top URLs**
- **Endpoint**: `/top-urls`
- **Method**: `GET`
- **Description**: Fetches the top N URLs. The results can be sorted either by the count of accesses (`count`) or by the latest submission time (`latest`).
- **Query Parameters**:
    - `sort_by`: Sorting criterion. Valid values are `"count"` or `"latest"`.
    - `get_n`: Number of top URLs to return.
- **Example Request**:
  ```bash
  curl "http://localhost:8080/top-urls?sort_by=count&get_n=50"
  ```
- **Response** (JSON):
  ```json
  [
    {
      "url": "http://example.com",
      "count": 50
    },
    {
      "url": "http://example2.com",
      "count": 40
    }
  ]
  ```

### 3. **Error Responses**
- **Invalid `sort_by`**: Returns `400 Bad Request` if an invalid value is provided for the `sort_by` parameter.
    - Example: `"sort_by": "invalid"`
- **Invalid `get_n`**: Returns `400 Bad Request` if `get_n` is not a valid integer.
    - Example: `"get_n": "not-a-number"`

## Batch Process

The application includes a **Batch Process** that runs periodically to collect and process the top URLs. It fetches the top 50 URLs from the store (by count), refetches them, updates their stats in the store and logs their stats. This process helps monitor URL activity and provides insights into the number of successes, failures, and the last download time for the top URLs.

### Configuration

The behavior of the batch process is controlled by the following parameters from the `config.yaml`:

- `worker_pool_size`: The number of concurrent workers used in processing the URLs.
- `num_of_batch_urls`: The number of top URLs to be collected and processed in each batch.
- `batch_interval_seconds`: The interval (in seconds) between each batch process execution.


## Internal Structure

### Data Model
URLs are stored in a linked list format to efficiently track and update URLs. Each URL has associated metadata::
- **Count**: Total number of submissions for the URL.
- **Successes**: Number of successful downloads for the URL.
- **Failures**: Number of failed download attempts.
- **Last Submitted**: Timestamp of the last submission.

The linked list structure allows for O(1) updates when a URL is added or modified.

## Setup

### Requirements
- Go 1.18+.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/url-submission-api.git
   cd url-submission-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

The server will start at `http://localhost:8080`.

## Configuration Overview

The application uses a YAML configuration file, `config.yaml`, to load various settings for the server and downloader. The config file is parsed into a Go struct, and the values are used to configure different aspects of the program's behavior.

### YAML Configuration Structure

The `config.yaml` file contains two main sections:

1. **server**: Configuration for the HTTP server
    - `port`: The port on which the HTTP server will listen for incoming requests. For example, `":8080"` will start the server on port 8080.

2. **downloader**: Configuration for the downloader's behavior
    - `worker_pool_size`: The number of concurrent worker goroutines to use in the downloader's worker pool. This controls how many URLs can be processed concurrently.
    - `num_of_batch_urls`: The number of URLs to process in each background batch process.
    - `batch_interval_seconds`: The interval, in seconds, between processing URL batches.

### Example Configuration

```yaml
server:
  port: ":8080"

downloader:
  worker_pool_size: 3
  num_of_batch_urls: 10
  batch_interval_seconds: 10
```

## Example Workflow

1. **Submit a URL**:
   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"url": "http://example.com"}' http://localhost:8080/submit-url
   ```

2. **Get Top URLs by Count**:
   ```bash
   curl "http://localhost:8080/top-urls?sort_by=count&get_n=5"
   ```

3. **Get Top URLs by Latest Submission**:
   ```bash
   curl "http://localhost:8080/top-urls?sort_by=latest&get_n=5"
   ```


# Potential enhancements
- Shut down workers when there's no tasks and spin them back up when necessary
- Use test data rather than executing actual gets
- Make two filter functions and make n not configurable to preallocate slice size and avoid reallocation
- Potentially shard the store (although I saw worse results as sorting didn't seem to have a massive overhead when benchmarking)
- How much do we care about accurate results? Could we batch sorting / processing to reduce overhead on fetching sorted lists
- Better worker pool & HTTP tests, most time spent on the store
- 