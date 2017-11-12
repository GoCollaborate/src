# GoCollaborate
[![Build Status](https://travis-ci.org/GoCollaborate/tests.svg?branch=master)](https://travis-ci.org/GoCollaborate/tests)
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
- [Source code](https://github.com/GoCollaborate/src)
- [Examples](https://github.com/GoCollaborate/examples)
- [Document](https://hastingsyoung.gitbooks.io/gocollaborateapi/content/)
- [GoCollaborate UI](https://github.com/HastingsYoung/Material-Dashboard-UI/tree/gocollaborate-ui)
- [Change Logs](https://github.com/GoCollaborate/src/blob/master/CHANGELOG.md)
- [Tests](https://github.com/GoCollaborate/tests)
## Quick Start
### Installation
```sh
go get -u github.com/GoCollaborate/src
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
	"github.com/GoCollaborate/src"
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
	"github.com/GoCollaborate/src/artifacts/task"
	"github.com/GoCollaborate/src/wrappers/taskHelper"
	"net/http"
)

func ExampleJobHandler(w http.ResponseWriter, r *http.Request) *task.Job {
	job := task.MakeJob()
	job.Tasks(&task.Task{task.SHORT,
		task.BASE, "exampleFunc",
		task.Collection{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4},
		task.Collection{0},
		task.NewTaskContext(struct{}{}), 0})
	job.Stacks("core.ExampleTask.Mapper", "core.ExampleTask.Reducer")
	return job
}

func ExampleFunc(source *task.Collection,
	result *task.Collection,
	context *task.TaskContext) bool {
	// deal with passed in request
	fmt.Println("Example Task Executed...")
	var total int
	// the function will calculate the sum of source data
	for _, n := range *source {
		total += n.(int)
	}
	result.Append(total)
	return true
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
