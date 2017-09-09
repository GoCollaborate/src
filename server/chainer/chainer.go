package chainer

import (
	"fmt"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/server/mapper"
	"github.com/GoCollaborate/server/reducer"
	"github.com/GoCollaborate/server/task"
	"sync"
	"time"
)

type ChainMapper interface {
	mapper.Mapper
	Append(m ...mapper.Mapper)
	Set(m mapper.Mapper, i int)
}

type BaseChainMapper struct {
	mappers []mapper.Mapper
}

func DefaultChainMapper(length int) *BaseChainMapper {
	return &BaseChainMapper{make([]mapper.Mapper, length)}
}

func (mp *BaseChainMapper) Append(m ...mapper.Mapper) {
	mp.mappers = append(mp.mappers, m...)
	return
}

func (mp *BaseChainMapper) Set(m mapper.Mapper, i int) {
	if i < len(mp.mappers) && i >= 0 {
		mp.mappers[i] = m
	}
}

func (mp *BaseChainMapper) Map(t *task.Task) (map[int64]*task.Task, error) {
	var (
		maps map[int64]*task.Task
		err  error
		init bool         = false
		lock sync.RWMutex = sync.RWMutex{}
	)
	for _, m := range mp.mappers {
		if m != nil {
			if !init {
				t.Context.Set("index", int64(0))
				maps, err = m.Map(t)
				init = true
				if err != nil {
					return maps, err
				}
			} else {
				cpmap := maps
				counter := 0
				l := len(maps)
				wait := make(chan bool)
				for k_m, mm := range cpmap {
					go func() {
						mm.Context.Set("index", k_m)
						_maps, err := m.Map(mm)

						if err != nil {
							return
						}

						lock.Lock()
						delete(maps, k_m)

						for k, value := range _maps {
							maps[k] = value
						}

						lock.Unlock()
						wait <- false
					}()
				}
				for {
					if counter == l {
						break
					}
					select {
					case <-wait:
						counter++
					case <-time.After(constants.DefaultMaxMappingTime):
						return maps, constants.ErrTimeout
					}
				}
			}
		}
	}
	fmt.Println(maps)
	return maps, nil
}

type PipelineMapper interface {
	mapper.Mapper
	Set(start mapper.Mapper, middle reducer.Reducer, end mapper.Mapper)
}

type maptuple struct {
	Start  mapper.Mapper
	Middle reducer.Reducer
	End    mapper.Mapper
}

func DefaultPipelineMapper() *BasePipelineMapper {
	return &BasePipelineMapper{nil}
}

type BasePipelineMapper struct {
	tuple *maptuple
}

func (p *BasePipelineMapper) Set(start mapper.Mapper, middle reducer.Reducer, end mapper.Mapper) {
	p.tuple = &maptuple{start, middle, end}
}

func (p *BasePipelineMapper) Map(t *task.Task) (map[int64]*task.Task, error) {
	var (
		maps map[int64]*task.Task
		err  error
	)
	maps, err = p.tuple.Start.Map(t)

	if err != nil {
		return maps, err
	}

	t = new(task.Task)

	err = p.tuple.Middle.Reduce(maps, t)

	if err != nil {
		return maps, err
	}

	maps, err = p.tuple.End.Map(t)

	return maps, nil
}
