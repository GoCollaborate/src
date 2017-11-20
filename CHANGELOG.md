# Change Logs
## Updates
**(Please note that no downward compability will be guaranteed before the formal release 1.0.0 )**
### 0.4.x
#### 0.4.0
- Update project distributed protocol to protobuf
- Update RPC framework to gRPC
- Extend automatic retry timeout interval
### 0.3.x
#### 0.3.2
- Add versioning for coordinator
- Support more load-balancing algorithms
- Optimize error responses
#### 0.3.1
- Add "methods" field to service response
- Change rate limiting response header: 422 -> 429 Too Many Requests
#### 0.3.0
- Add Travis CI support
- The previous repository was deprecated, now please refer to [https://github.com/GoCollaborate/src](https://github.com/GoCollaborate/src)
### 0.2.x
#### 0.2.9
- Optimize UI
- Add comments
#### 0.2.8
- Update entry API -> []task.Countable -> task.Collection
- Provide built-in functions for task.Collection
- Update examples
- Update entry API -> function return type changed to bool
- Add pagination to UI
#### 0.2.7
- Add seeds to Gossip protocol
- Support reading streams of CSV file from various soruces
- Update examples and test scripts
#### 0.2.6
- Update service register heartbeats
- Add GC for expired services
- Add TaskHelper utility
- Create API test packs
#### 0.2.5
- Add analytics support
- Optimize UI
#### 0.2.4
- Optimize performance
- Optimize UI
#### 0.2.3
- Rewrite the project structrue
- Support full gossipping
- Support rate limiting
- Add utils to taskHelper
#### 0.2.2
- Integrate with GoCollaborate UI (http://localhost:8080)
- Refine web routes structrue
#### 0.2.1
- Refactor Task API
- Refine communication structs
- Support generic executors
- Rewrite examples
#### 0.2.0
- Implement Gossip Protocol
- Rename struct and functions
- Refine code structure
### 0.1.x
#### 0.1.9
- Refactor package dependencies
- Add Job, Stage literals
- Support bulk-execution of tasks
- Update example documents
#### 0.1.8
- Refactor API entry in Coordinator mode
- Add documents for Coordinator mode
- Add command line argument `numwks` to specify number of workers per master program
#### 0.1.7
- Support chainning of mappers
- Support mapper & reducer pipelines
- Refactor example project structrue
#### 0.1.6
- Support task-specific mapper and reducer
- Update example documents
- Repair bugs in coordinator mode
#### 0.1.5
- Refine logging format
- Support automatic garbage collection for hash functions