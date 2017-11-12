package stats

import (
	"github.com/GoCollaborate/src/constants"
	"sync"
	"time"
)

var (
	once      sync.Once
	singleton *StatsManager
)

type StatsManager struct {
	stats map[string]*AbstractArray
	// channel of hits
	StatsChan map[string]chan Hit
	// accumulative hits
	StatsAcc      map[string]*[]Hit
	StatsPolicies map[string]*AbsPolicy
}

// Hit is a single hit of stat record
type Hit struct {
	// stat key
	Key string `json:"key,omitempty"`
	// stat value
	Val interface{} `json:"val"`
}

type AbstractArray struct {
	// the abstract array
	Array []Abstract `json:"array"`
	// the policy name to generate such abstract
	Policy string `json:"policy"`
}

// Abstract is the abstract of Hits over a particular period of time
// An abstract is usually set to the mean, sum or residue of Hits
type Abstract struct {
	// stat key
	Key string `json:"key,omitempty"`
	// stat value
	Val interface{} `json:"val"`
	// the middle value timestamp of the stat object in epoch second
	MidTs int64 `json:"midts"`
	// the length of stat object in second
	Len float64 `json:"len,omitemtpy"`
}

type AbsPolicy struct {
	Name  string
	Funct func(vals ...Hit) interface{}
}

func GetStatsInstance() *StatsManager {
	once.Do(func() {
		var (
			pol = AbsPolicySumOfInt()
		)

		singleton = &StatsManager{
			map[string]*AbstractArray{},
			map[string]chan Hit{},
			map[string]*[]Hit{},
			map[string]*AbsPolicy{constants.StatsPolicySumOfInt: pol},
		}

		singleton.Observe("tasks")
		singleton.Observe("hits")

		// flush hits from stats channel to hits array
		go func() {
			for {
				singleton.flush()
				<-time.After(constants.DefaultStatFlushInterval)
			}
		}()

		// generate abstract from hits array
		go func() {
			for {
				singleton.abstract()
				<-time.After(constants.DefaultStatAbstractInterval)
			}
		}()
	})
	return singleton
}

// Record a hit
func (sm *StatsManager) Record(t string, v interface{}, k ...string) error {
	if ch, ok := sm.StatsChan[t]; ok {
		if len(k) > 0 {
			// asyn call will reduce performance issue regardless of flush interval
			go func() { ch <- Hit{Key: k[0], Val: v} }()
			return nil
		}
		// asyn call will reduce performance issue regardless of flush interval
		go func() { ch <- Hit{Val: v} }()
		return nil
	}
	return constants.ErrStatTypeNotFound
}

// Specify the custom route to observe
func (sm *StatsManager) Observe(route string) {
	var (
		arr  = DefaultAbstractArray()
		hits = make(chan Hit)
		acc  = &[]Hit{}
	)
	sm.stats[route] = arr
	sm.StatsChan[route] = hits
	sm.StatsAcc[route] = acc
}

// Return a map of routes - AbstractArray
func (sm *StatsManager) Stats() map[string]*AbstractArray {
	return sm.stats
}

// Flush hits from observed channels into hits array
func (sm *StatsManager) flush() {
	for arrk, _ := range sm.stats {
		var (
			vs *[]Hit
			ch chan Hit
			ok = false
		)

		if ch, ok = sm.StatsChan[arrk]; !ok {
			// no channel match against the abstract, skip current loop
			continue
		}

		if vs, ok = sm.StatsAcc[arrk]; !ok {
			// no hit array match against the abstract, skip current loop
			continue
		}

		select {
		case v := <-ch:
			*vs = append(*vs, v)
		default:
			ok = false
		}
	}
}

// Generate abstract from hits array
func (sm *StatsManager) abstract() {
	for arrk, arr := range sm.stats {
		var (
			vs  *[]Hit
			pol *AbsPolicy
			ok  = false
		)

		if pol, ok = sm.StatsPolicies[arr.Policy]; !ok {
			// no policy match against the abstract
			continue
		}

		if vs, ok = sm.StatsAcc[arrk]; !ok {
			// no hit array match against the abstract
			continue
		}

		// extended array
		ext := *vs
		*vs = []Hit{}

		// no extension
		if len(ext) < 1 {
			continue
		}

		// append extension to abstract array
		arr.Array = append(arr.Array,
			Abstract{
				Val:   pol.Funct(ext...),
				Len:   constants.DefaultStatAbstractInterval.Seconds(),
				MidTs: time.Now().Add(-constants.DefaultStatAbstractInterval / 2).Unix()})
	}
}

func DefaultAbstractArray() *AbstractArray {
	return &AbstractArray{[]Abstract{}, constants.StatsPolicySumOfInt}
}

func AbsPolicySumOfInt() *AbsPolicy {
	return &AbsPolicy{constants.StatsPolicySumOfInt, func(vals ...Hit) interface{} {
		sum := 0
		for _, v := range vals {
			sum += v.Val.(int)
		}
		return sum
	}}
}
