// Gian Lucas Martín Chamorro
package main

import (
	"fmt"
	"runtime"
	"time"
)

const (
	Procs = 4

	H = 10

	Bears    = 1
	EatCount = 5
	TimeEat  = 1

	Bees         = 5
	ProduceCount = 10
	TimeProduce  = 2
)

type Empty struct{}

func eat() {
	time.Sleep(TimeEat * time.Millisecond)
}

func produce() {
	time.Sleep(TimeProduce * time.Millisecond)
}

func bear(done, HoneyJar, BearsBarrier, BeesBarrier chan Empty) {
	for i := 0; i < EatCount; i++ {
		<-BearsBarrier
		eat()
		for j := 0; j < H; j++ {
			<-HoneyJar
		}
		fmt.Printf("\t\tBear: Honey jar is now empty\n")
		BeesBarrier <- Empty{}
	}
	done <- Empty{}
}

func bee(id int, done chan Empty, HoneyJar chan Empty, BearsBarrier chan Empty, BeesBarrier chan Empty) {
	for i := 0; i < ProduceCount; i++ {
		<-BeesBarrier
		produce()
		HoneyJar <- Empty{}
		if len(HoneyJar) < cap(HoneyJar) {
			fmt.Printf("Bee %d: Portion number %d is on the jar now\n", id, len(HoneyJar))
			BeesBarrier <- Empty{}
		} else {
			fmt.Printf("Bee %d: Honey jar is now full\n", id)
			BearsBarrier <- Empty{}
		}
	}
	done <- Empty{}
}

func main() {
	runtime.GOMAXPROCS(Procs)
	done := make(chan Empty, 1)
	HoneyJar := make(chan Empty, H)
	BeesBarrier := make(chan Empty, 1)
	BearsBarrier := make(chan Empty, 1)

	BeesBarrier <- Empty{}
	//Bear
	for i := 0; i < Bears; i++ {
		go bear(done, HoneyJar, BearsBarrier, BeesBarrier)
	}

	//Bees
	for i := 0; i < Bees; i++ {
		go bee(i, done, HoneyJar, BearsBarrier, BeesBarrier)
	}

	for i := 0; i < Bees+Bears; i++ {
		<-done
	}

}
