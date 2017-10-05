package task

type Stage struct {
	previous *Stage
	next     *Stage
	TaskSet  map[int]*Task
}

func (s *Stage) Prev() *Stage {
	return s.previous
}

func (s *Stage) Next() *Stage {
	return s.next
}

func MakeStage(prev *Stage, next *Stage, tset ...map[int]*Task) *Stage {
	if len(tset) > 0 {
		var maps map[int]*Task = make(map[int]*Task)
		for _, val := range tset {
			for k, v := range val {
				maps[k] = v
			}
		}
		return &Stage{prev, next, maps}
	}
	return &Stage{prev, next, map[int]*Task{}}
}
