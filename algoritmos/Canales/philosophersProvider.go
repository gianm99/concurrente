package main

import (
	"fmt"
	"runtime"
	"time"
)

const (
	procs         = 4
	cphilosophers = 4
	eatCount      = 5

	thinking = iota
	hungry
	eating
)

type empty struct{}
type request struct {
	c      chan empty
	id     int
	status int
}

func main() {
	runtime.GOMAXPROCS(procs)
	done := make(chan empty, 1)
	providerChan := make(chan request)
	go provider(providerChan)
	for i := 0; i < cphilosophers; i++ {
		go philosopher(i, done, providerChan)
	}
	for i := 0; i < cphilosophers; i++ {
		<-done
	}
	fmt.Println("End")
}

func philosopher(id int, done chan empty, provider chan request) {
	myChan := make(chan empty)
	think := func() {
		time.Sleep(50 * time.Millisecond)
	}
	eat := func() {
		time.Sleep(50 * time.Millisecond)
	}
	pick := func() {
		provider <- request{id: id, c: myChan, status: hungry}
		<-myChan
	}
	release := func() {
		provider <- request{id: id, c: myChan, status: thinking}
	}
	for i := 0; i < eatCount; i++ {
		think()
		pick()
		eat()
		release()
	}
	done <- empty{}
}

func provider(channel chan request) {
	var philosophers [cphilosophers]request

	right := func(i int) int {
		return (i + 1) % cphilosophers
	}

	left := func(i int) int {
		return (i + cphilosophers - 1) % cphilosophers
	}

	canEat := func(i int) {
		r := right(i)
		l := left(i)
		if philosophers[i].status == hungry &&
			philosophers[l].status != eating &&
			philosophers[r].status != eating {
			philosophers[i].status = eating
			philosophers[i].c <- empty{}
		}
	}
	for i := range philosophers {
		philosophers[i].status = thinking
	}
	for {
		m := <-channel
		philosophers[m.id] = m
		switch m.status {
		case hungry:
			canEat(m.id)
		case thinking:
			canEat(left(m.id))
			canEat(right(m.id))
		}
	}
}
