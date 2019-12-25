package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type transaccion struct {
	permiso bool
	balance int32
}

// Pasa una transacción a una secuencia de bytes
func transtobyte(trans transaccion) []byte {
	b := make([]byte, 5)
	b[0] = byte(Btoi(trans.permiso))
	binary.LittleEndian.PutUint32(b[1:], uint32(trans.balance))
	return b
}

// Btoi pasa un bool a int
func Btoi(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
