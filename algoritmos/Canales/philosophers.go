package main

import (
	"fmt"
	"runtime"
	"time"
)

const (
	procs        = 4
	philosophers = 4
	eatCount     = 100
)

type empty struct{}
type fork chan empty

type tforks [philosophers]fork

func think(id int) {
	fmt.Printf("%d start sleep\n", id)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%d end sleep\n", id)
}

func eat(id int) {
	fmt.Printf("%d start eat\n", id)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("%d end eat\n", id)
}

func right(i int) int {
	return (i + 1) % philosophers
}

func pick(id int, forks tforks) {
	if id < right(id) {
		<-forks[id]
		<-forks[right(id)]
	} else {
		<-forks[right(id)]
		<-forks[id]
	}
}

func release(id int, forks tforks) {
	forks[id] <- empty{}
	forks[right(id)] <- empty{}
}

func philosopher(id int, done chan empty, forks tforks) {
	for i := 0; i < eatCount; i++ {
		think(id)
		pick(id, forks)
		eat(id)
		release(id, forks)
	}
	done <- empty{}
}

func main() {
	runtime.GOMAXPROCS(procs)
	done := make(chan empty, 1)
	var forks tforks
	for i := range forks {
		forks[i] = make(chan empty, 1)
		forks[i] <- empty{}
	}
	for i := 0; i < philosophers; i++ {
		go philosopher(i, done, forks)
	}
	for i := 0; i < philosophers; i++ {
		<-done
	}
	fmt.Println("End")
}
