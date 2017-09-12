# GoCollaborate
## What is GoCollaborate?
GoCollaborate is an universal framework for distributed services management that you can easily program with, build extension on, and on top of which you can create your own high performance distributed services like a breeze.
## The Idea Behind
GoCollaborate absorbs the best practice experience from popular distributed services frameworks like[✨Hadoop](https://hadoop.apache.org/), [✨ZooKeeper](https://zookeeper.apache.org/), [✨Dubbo](http://dubbo.io/) and [✨Kite](https://github.com/koding/kite) that helps to ideally resolve the communication and collaboration issues among multiple isolated peer servers.
## Am I Free to Use GoCollaborate?
Yes! Please check out the terms of the BSD License.
## Contribution
This project is currently under development, please feel free to fork it and report issues!
## Documents (In Construction...)
![alt text](https://github.com/HastingsYoung/GoCollaborate/raw/master/home.png "Docs Home Page")
## Relative Links
- [Source code](https://github.com/HastingsYoung/GoCollaborate)
- [Examples](https://github.com/HastingsYoung/GoCollaborateExamples)
## Updates (Please note that no downward compability will be guaranteed before the formal release 1.0.0 )
### 0.1.8
- Refactor API entry in Coordinator Mode
- Add documents for Coordinator mode
- Add command line argument `numwks` to specify number of workers per master program
### 0.1.7
- Support chainning of mappers
- Support mapper & reducer pipelines
- Refactor example project structrue
### 0.1.6
- Support task-specific mapper and reducer
- Update example documents
- Repair bugs in coordinator mode
### 0.1.5
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
touch contact.json
touch main.go
cd ./core
touch simpleTaskHandler.go
```
The project structure now looks something like this:
```
[Your_Project_Name]
┬
├ [core]
	┬
	└ simpleTaskHandler.go
├ contact.json
└ main.go
```
Configure file `contact.json`:
```js
{
	"agents": [{
		"ip": "localhost",
		"port": 57851,
		"alive": true
	}, {
		"ip": "localhost",
		"port": 57852,
		"alive": true
	}],
	"local": {
		"ip": "localhost",
		"port": 57851,
		"alive": true
	},
	"coordinator": {
		"ip": "",
		"port": 0,
		"alive": false
	},
	"timestamp": 1504182648
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
	collaborate.Set("Mapper", mp, "core.ExampleTaskHandler.Mapper")
	collaborate.Set("Reducer", rd, "core.ExampleTaskHandler.Reducer")
	collaborate.Set("Shared", []string{"GET", "POST"}, core.ExampleTaskHandler)
	collaborate.Run()
}

```
### Map-Reduce
```go
package core

import (
	"fmt"
	"github.com/GoCollaborate/server/task"
	"net/http"
)

func ExampleTaskHandler(w http.ResponseWriter, r *http.Request) task.Task {
	return task.Task{task.PERMANENT,
		task.BASE, "exampleFunc", []task.Countable{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4},
		[]task.Countable{0},
		task.NewTaskContext(struct{}{}), "core.ExampleTaskHandler.Mapper", "core.ExampleTaskHandler.Reducer"}
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

func (m *SimpleMapper) Map(t *task.Task) (map[int64]*task.Task, error) {
	maps := make(map[int64]*task.Task)
	s1 := t.Source[:4]
	s2 := t.Source[4:8]
	s3 := t.Source[8:]
	s4 := t.Result
	s5 := t.Result
	s6 := t.Result
	maps[int64(0)] = &task.Task{t.Type, t.Priority, t.Consumable, s1, s4, t.Context, t.Mapper, t.Reducer}
	maps[int64(1)] = &task.Task{t.Type, t.Priority, t.Consumable, s2, s5, t.Context, t.Mapper, t.Reducer}
	maps[int64(2)] = &task.Task{t.Type, t.Priority, t.Consumable, s3, s6, t.Context, t.Mapper, t.Reducer}
	return maps, nil
}

type SimpleReducer int

func (r *SimpleReducer) Reduce(source map[int64]*task.Task, result *task.Task) error {
	rs := *result
	var sum int
	for _, s := range source {
		for _, r := range (*s).Result {
			sum += r.(int)
		}
	}
	rs.Result[0] = sum
	fmt.Printf("The sum of numbers is: %v \n", sum)
	fmt.Printf("The task result set is: %v", rs)
	return nil
}

```
### Run
Here we create the entry file and a simple implementation of map-reduce interface, and next we will run with std arguments:
```sh
go run main.go -svrmode=clbt
```
The task is now up and running at:
```
http://localhost:8080/core/ExampleTaskHandler
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
3. Edit local ip address in `contact.json`:
```js
{
	"agents": [{
		"ip": "localhost",
		"port": 57851,
		"alive": true
	}, {
		"ip": "localhost",
		"port": 57852,
		"alive": true
	}],
	"local": {
		"ip": "localhost",
		"port": 57852,
		"alive": true
	},
	"coordinator": {
		"ip": "",
		"port": 0,
		"alive": false
	},
	"timestamp": 1504182648
}
```
4. Run the copied project, don't forget to change your port number if you are running locally:
```sh
go run main.go -svrmode=clbt -port=8081
```
5. Now the distributed servers are available at:
```
http://localhost:8080/core/ExampleTaskHandler
// and 
http://localhost:8081/core/ExampleTaskHandler
```

## Acknowledgement
- [fatih/color](https://github.com/fatih/color)
- [mattn/go-colorable](https://github.com/mattn/go-colorable)
- [mattn/go-isatty](https://github.com/mattn/go-isatty)
- [gorilla/mux](https://github.com/gorilla/mux)
