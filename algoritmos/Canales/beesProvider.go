package main

import (
	"fmt"
	"time"
)

const (
	nbees   = 4
	jarsize = 10
)

type empty struct{}
type request struct {
	id int
	c  chan empty
}

func main() {
	doneEating := make(chan empty)
	wakeUp := make(chan empty)
	jarReq := make(chan request)
	forever := make(chan empty)
	go jar(jarReq, wakeUp, doneEating)
	go bear(wakeUp, doneEating)
	for i := 0; i < nbees; i++ {
		go bee(i, jarReq)
	}
	<-forever
}

func bear(wakeUp chan empty, doneEating chan empty) {
	fmt.Printf("Bear begins\n")
	for {
		<-wakeUp
		fmt.Println("Bear starts to eat")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Bear finishes eating")
		doneEating <- empty{}
	}
}

func bee(id int, jarReq chan request) {
	myChan := make(chan empty)
	fmt.Printf("Bee %d begins\n", id)
	for {
		time.Sleep(50 * time.Millisecond)
		jarReq <- request{id, myChan}
		<-myChan
	}
}

func jar(channel chan request, wakeUp chan empty, doneEating chan empty) {
	fmt.Println("Honey jar (Provider) begins")
	portions := 0
	for {
		m := <-channel
		portions++
		if portions == jarsize {
			wakeUp <- empty{}
			fmt.Printf("Bee %d wakes the bear\n", m.id)
			<-doneEating
			portions = 0
		}
		m.c <- empty{}
	}
}
