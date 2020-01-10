package Canales

type empty struct{}
type mutex chan empty

func newMutex() mutex {
	m := make(mutex, 1)
	m <- empty{}
	return m
}

func (m mutex) lock() {
	<-m
}

func (m mutex) unlock() {
	m <- empty{}
}
