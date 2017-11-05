# GoCollaborate
![alt text](https://github.com/HastingsYoung/GoCollaborate/raw/master/home.png "Docs Home Page")
## What is GoCollaborate?
GoCollaborate is an universal framework for stream computing and distributed services management that you can easily program with, build extension on, and on top of which you can create your own high performance distributed applications.
### The Idea Behind
GoCollaborate absorbs the best practice experience and improves from the popular distributed computing frameworks including[✨Hadoop](https://hadoop.apache.org/), [✨Spark](https://spark.apache.org/), [✨ZooKeeper](https://zookeeper.apache.org/), [✨Dubbo](http://dubbo.io/) and [✨Kite](https://github.com/koding/kite) that helps to ideally provision the computability for large scale data sets with an easy-to-launch setups.
### Am I Free to Use GoCollaborate?
Yes! Please check out the terms of the BSD License.
### Contribution
This project is currently under development, please feel free to fork it and report issues!

Please check out most recent [API](https://hastingsyoung.gitbooks.io/gocollaborateapi/content/) document for more information.

### Relative Links
- [Source code](https://github.com/HastingsYoung/GoCollaborate)
- [Examples](https://github.com/HastingsYoung/GoCollaborateExamples)
- [Document](https://hastingsyoung.gitbooks.io/gocollaborateapi/content/)
- [GoCollaborate UI](https://github.com/HastingsYoung/Material-Dashboard-UI/tree/gocollaborate-ui)
## Updates
**(Please note that no downward compability will be guaranteed before the formal release 1.0.0 )**
### 0.2.x
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
## Quick Start
### Installation
```sh
go get -u github.com/GoCollaborate
```
### Create Project
```sh
mkdir Your_Project_Name
cd Your_Project_Name
mkdir core
touch case.json
touch main.go
cd ./core
touch example.go
```
The project structure now looks something like this:
```
[Your_Project_Name]
┬
├ [core]
	┬
	└ example.go
├ case.json
└ main.go
```
Configure file `case.json`:
```js
{
	"caseid": "GoCollaborateStandardCase",
	"cards": {
		"localhost:57851": {
			"ip": "localhost",
			"port": 57851,
			"alive": false,
			"seed": false
		},
		"localhost:57852": {
			"ip": "localhost",
			"port": 57852,
			"alive": true,
			"seed": true
		}
	},
	"timestamp": 1508619931,
	"local": {
		"ip": "localhost",
		"port": 57852,
		"alive": true,
		"seed": true
	},
	"coordinator": {
		"ip": "localhost",
		"port": 0,
		"alive": true,
		"seed": false
	}
}
```
### Entry
```go
package main

import (
	"./core"
	"github.com/GoCollaborate"
)

func main() {
	mp := new(core.SimpleMapper)
	rd := new(core.SimpleReducer)
	collaborate.Set("Function", core.ExampleFunc, "exampleFunc")
	collaborate.Set("Mapper", mp, "core.ExampleTask.Mapper")
	collaborate.Set("Reducer", rd, "core.ExampleTask.Reducer")
	collaborate.Set("Shared", []string{"GET", "POST"}, core.ExampleJobHandler)
	collaborate.Run()
}

```
### Map-Reduce
```go
package core

import (
	"fmt"
	"github.com/GoCollaborate/artifacts/task"
	"github.com/GoCollaborate/wrappers/taskHelper"
	"net/http"
)

func ExampleJobHandler(w http.ResponseWriter, r *http.Request) *task.Job {
	job := task.MakeJob()
	job.Tasks(&task.Task{task.SHORT,
		task.BASE, "exampleFunc",
		[]task.Countable{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4},
		[]task.Countable{0},
		task.NewTaskContext(struct{}{}), 0})
	job.Stacks("core.ExampleTask.Mapper", "core.ExampleTask.Reducer")
	return job
}

func ExampleFunc(source *[]task.Countable,
	result *[]task.Countable,
	context *task.TaskContext) chan bool {
	out := make(chan bool)
	// deal with passed in request
	go func() {
		fmt.Println("Example Task Executed...")
		var total int
		for _, n := range *source {
			total += n.(int)
		}
		*result = append(*result, total)
		out <- true
	}()
	return out
}

type SimpleMapper int

func (m *SimpleMapper) Map(inmaps map[int]*task.Task) (map[int]*task.Task, error) {
	// slice the data source of the map into 3 separate segments
	return taskHelper.Slice(inmaps, 3), nil
}

type SimpleReducer int

func (r *SimpleReducer) Reduce(maps map[int]*task.Task) (map[int]*task.Task, error) {
	var sum int
	for _, s := range maps {
		for _, r := range (*s).Result {
			sum += r.(int)
		}
	}
	fmt.Printf("The sum of numbers is: %v \n", sum)
	fmt.Printf("The task set is: %v", maps)
	return maps, nil
}

```
### Run
Here we create the entry file and a simple implementation of map-reduce interface, and next we will run with std arguments:
```sh
go run main.go -mode=clbt
```
The task is now up and running at:
```
http://localhost:8080/core/ExampleJobHandler
```
### Collaborate
1. Copy your project directory:
```sh
cp Your_Project_Name Your_Project_Name_Copy
```
2. Enter the copied project:
```sh
cd Your_Project_Name_Copy
```
3. Edit local ip address in `case.json`:
```js
{
	"caseid": "GoCollaborateStandardCase",
	"cards": {
		"localhost:57852": {
			"ip": "localhost",
			"port": 57852,
			"alive": true,
			"seed": true
		}
	},
	"timestamp": 1508619931,
	"local": {
		"ip": "localhost",
		"port": 57851,
		"alive": true,
		"seed": false
	},
	"coordinator": {
		"ip": "localhost",
		"port": 0,
		"alive": true,
		"seed": false
	}
}
```
4. Run the copied project, don't forget to change your port number if you are running locally:
```sh
go run main.go -mode=clbt -port=8081
```
5. Now the distributed servers are available at:
```
http://localhost:8080/core/ExampleJobHandler
// and 
http://localhost:8081/core/ExampleJobHandler
```
6. Alternatively, access the GoCollaborate UI for more infomation:
```
http://localhost:8080
```

## Acknowledgement
- [fatih/color](https://github.com/fatih/color)
- [mattn/go-colorable](https://github.com/mattn/go-colorable)
- [mattn/go-isatty](https://github.com/mattn/go-isatty)
- [gorilla/mux](https://github.com/gorilla/mux)
- [satori/go.uuid](https://github.com/satori/go.uuid)
- [radioinmyhead/csv](https://github.com/radioinmyhead/csv)
