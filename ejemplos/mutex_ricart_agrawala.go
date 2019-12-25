/* The Ricart-Agrwala distributed mutual exclusion algorithm */

package main

import (
    "fmt"
    "runtime"
    "sync"
)

type Empty struct{}

type Message struct {
    source int
    number int
}

const (
    Procs    = 4
    MaxCount = 10000000
    Nodes    = 4
)

var counter = 0

func node(id, counts int, done chan Empty, requests, replies [Nodes]chan Message) {
    myNumber := 0
    deferred := make(chan int, Nodes)
    highestNum := 0
    requestCS := false
    mutex := new(sync.Mutex)//Semafor binari de la llibreria sync

    /* This is the asynchronous thread to receive requests from othe nodes*/
    receiver := func() {
        for {
            m := <-requests[id] //rep un missatge i compara amb el torn del node
            mutex.Lock()
            if m.number > highestNum { // actualitza el nombre mes gran
                highestNum = m.number  // no hi es a l'algorisme general
            }
            if !requestCS || (m.number < myNumber ||
                (m.number == myNumber && m.source < id)) { // prodra passa a la SC
                mutex.Unlock()
                replies[m.source] <- Message{source: id}   // envia permis
            } else {                                       // no pot passar
                deferred <- m.source                       // a la cua de diferits sense missatge
                mutex.Unlock()
            }
        }
    }

    // Launch the receiver
    go receiver()

    lock := func() {
        mutex.Lock()
        requestCS = true
        myNumber = highestNum + 1 // Assignacio de torn SC
        mutex.Unlock()
        for i := range requests {
            if i == id {
                continue
            }
            //Ennvia peticio a tots els altres
            requests[i] <- Message{source: id, number: myNumber}
        }
        for i := 0; i < Nodes-1; i++ {
            <-replies[id]//Espera rèplica de tots els altres
        }
    }

    unlock := func() {
        requestCS = false
        mutex.Lock()
        n := len(deferred)      // longitud de la cua
        mutex.Unlock()
        for i := 0; i < n; i++ {// Tos els diferits
            src := <-deferred   // el lleva de la cua
            replies[src] <- Message{source: id}
        }
    }

    for i := 0; i < counts; i++ {
        lock()
        counter++
        unlock()
    }

    fmt.Printf("End %d counter: %d\n", id, counter)
    done <- Empty{}
}

func main() {
    runtime.GOMAXPROCS(Procs)
    done := make(chan Empty, 1)
    // Arrays de 4 (Nodes) missatges (nom + num)
    var requests, replies [Nodes]chan Message

    for i := range replies {
        // peticions d'acces
        requests[i] = make(chan Message)
        // contestacio a les peticions
        replies[i] = make(chan Message)
    }

    for i := 0; i < Nodes; i++ {
        // node, identificació, el que ha de comptar, canal per acabar
        // array de peticions i de rèpliques
        go node(i, MaxCount/Nodes, done, requests, replies)
    }

    for i := 0; i < Nodes; i++ {
        <-done
    }

    fmt.Printf("Counter value: %d Expected: %d\n", counter, MaxCount)
}
