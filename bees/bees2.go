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
	TimeEat  = 3

	Bees         = 5
	ProduceCount = 10
	TimeProduce  = 1
)

type Empty struct{}

type Info struct {
	bee     string
	portion int
}

func eat() {
	time.Sleep(TimeEat * time.Millisecond)
}

func produce() {
	time.Sleep(TimeProduce * time.Millisecond)
}

func bear(done chan Empty, HoneyJar chan int, WakeUp chan Info) {

	for i := 0; i < EatCount; i++ {
		info := <-WakeUp
		fmt.Printf("\t\tBear: %s produced portion %d and woke me up\n", info.bee, info.portion)
		eat()
		fmt.Printf("\t\tBear: Honey jar is now empty\n")
		HoneyJar <- 0
	}
	done <- Empty{}
}

func bee(id int, done chan Empty, HoneyJar chan int, WakeUp chan Info) {
	for i := 0; i < ProduceCount; i++ {
		Honey := <-HoneyJar
		produce()
		Honey++
		if Honey < H {
			fmt.Printf("Bee %d: Portion number %d is on the jar now\n", id, Honey)
			HoneyJar <- Honey
		} else {
			fmt.Printf("Bee %d: Honey jar is now full\n", id)
			WakeUp <- Info{bee: fmt.Sprintf("Bee %d", id), portion: Honey}
		}
	}
	done <- Empty{}
}

func main() {
	runtime.GOMAXPROCS(Procs)
	done := make(chan Empty, 1)
	HoneyJar := make(chan int, 1)
	WakeUp := make(chan Info, 1)

	HoneyJar <- 0
	//Bear
	for i := 0; i < Bears; i++ {
		go bear(done, HoneyJar, WakeUp)
	}

	//Bees
	for i := 0; i < Bees; i++ {
		go bee(i, done, HoneyJar, WakeUp)
	}

	for i := 0; i < Bees+Bears; i++ {
		<-done
	}

}
