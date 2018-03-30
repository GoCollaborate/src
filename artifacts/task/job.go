package task

import (
	"github.com/GoCollaborate/src/constants"
	"github.com/satori/go.uuid"
	"sort"
)

type Job struct {
	jid    string
	stage  *Stage
	front  *Stage
	back   *Stage
	length int
	stacks [][]string
}

func MakeJob(s ...*Stage) *Job {
	if len(s) > 0 {
		return &Job{uuid.NewV4().String(), s[0], s[0], s[0], 1, [][]string{}}
	}
	return &Job{uuid.NewV4().String(), nil, nil, nil, 0, [][]string{}}
}

func (j *Job) Id() string {
	return j.jid
}

func (j *Job) Len() int {
	return j.length
}

func (j *Job) Back() *Stage {
	return j.back
}

func (j *Job) Curr() *Stage {
	return j.stage
}

func (j *Job) Front() *Stage {
	return j.front
}

func (j *Job) InsertBefore(bef *Stage, curr *Stage) *Stage {
	if curr == nil {
		return j.Init(bef)
	}
	if curr == j.front {
		return j.PushFront(bef)
	}

	bef.previous = curr.previous
	bef.next = curr
	if curr.previous != nil {
		curr.previous.next = bef
	}

	curr.previous = bef
	j.length++
	return bef
}

func (j *Job) InsertAfter(aft *Stage, curr *Stage) *Stage {
	if curr == nil {
		return j.Init(aft)
	}
	if curr == j.back {
		return j.PushBack(aft)
	}

	aft.next = curr.next
	aft.previous = curr
	if curr.next != nil {
		curr.next.previous = aft
	}

	curr.next = aft
	j.length++
	return aft
}

func (j *Job) Init(s *Stage) *Stage {
	if j.Len() == 0 {
		j.front = s
		j.stage = s
		j.back = s
		j.length++
		return s
	}
	return j.Curr()
}

func (j *Job) PushBack(back *Stage) *Stage {
	if j.Len() == 0 {
		return j.Init(back)
	}
	if j.back != nil {
		j.back.next = back
		back.previous = j.back
	}
	j.back = back
	j.length++
	return back
}

func (j *Job) PushFront(front *Stage) *Stage {
	if j.Len() == 0 {
		return j.Init(front)
	}
	if j.front != nil {
		j.front.previous = front
		front.next = j.front
	}
	j.front = front
	j.length++
	return front
}

func (j *Job) Exes(i int) ([]string, error) {
	if i > len(j.stacks) {
		return []string{}, constants.ERR_EXECUTOR_STACK_LENGTH_INCONSISTENT
	}
	return j.stacks[i], nil
}

func (j *Job) Stacks(stacks ...string) *Job {
	j.stacks = append(j.stacks, stacks)
	return j
}

func (j *Job) Tasks(tsks ...*Task) {
	if len(tsks) < 1 {
		return
	}
	sort.SliceStable(tsks, func(i, j int) bool {
		return tsks[i].Stage < tsks[j].Stage
	})
	var (
		stage int    = tsks[0].Stage
		temp  *Stage = MakeStage(nil, nil)
	)
	for i, s := range tsks {
		if stage == s.Stage {
			temp.TaskSet[i] = s
		} else {
			j.PushBack(temp)
			temp = MakeStage(nil, nil)
			temp.TaskSet[i] = s
			stage++
		}
	}
	j.PushBack(temp)
}
