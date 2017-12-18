package task

type Background chan *Job

func NewBackground() *Background {
	bg := make(Background)
	return &bg
}

func (bg *Background) Done() *Job {
	return <-*bg
}

func (bg *Background) Mount(job *Job) {
	go func() {
		*bg <- job
	}()
}

func (bg *Background) Close() {
	close(*bg)
}
