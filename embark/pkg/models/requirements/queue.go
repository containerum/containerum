package requirements

type Queue chan Dependency

func NewQueue(size int) Queue {
	return make(Queue, size)
}

func (queue Queue) Push(deps ...Dependency) {
	go func() {
		for _, dep := range deps {
			dep := dep
			queue <- dep
		}
	}()
}

func (queue Queue) Pop() *Dependency {
	if len(queue) > 0 {
		var dep = <-queue
		return &dep
	}
	return nil
}
