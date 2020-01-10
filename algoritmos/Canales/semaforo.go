package Canales

type empty struct{}
type semaphore struct {
	value chan int
	queue chan empty
}

func newSemaphore(value int) semaphore {
	var s semaphore
	if value < 0 {
		value = 0
	}
	s.value = make(chan int, 1)
	s.queue = make(chan empty, 1)
	s.value <- value
	return s
}

func (s semaphore) wait() {
	v := <-s.value
	v--
	if v < 0 {
		<-s.queue
	}
}

func (s semaphore) signal() {
	v := <-s.value
	v++
	s.value <- v
	if v <= 0 {
		s.queue <- empty{}
	}
}
