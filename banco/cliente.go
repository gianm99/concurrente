package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

const maxOp = 8

type transaccion struct {
	permiso bool
	balance int32
}

// Pasa la secuencia de bytes a una transacción
func bytetotrans(arr []byte) transaccion {
	permiso := int8(arr[0]) != 0
	balance := int32(binary.LittleEndian.Uint32(arr[1:]))
	return transaccion{permiso, balance}
}

// operacionRPC intenta realizar una operación bancaria
func operacionRPC(cantidad int32) (trans transaccion, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fallo al conectar con RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Fallo al abrir un canal")
	defer ch.Close()

	// Cola para respuestas del banquero
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Fallo al declarar una cola")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Fallo al registrar un consumidor")

	corrId := randomString(32)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(cantidad))
	// Mensaje para hacer la operacion
	err = ch.Publish(
		"",            // exchange
		"operaciones", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          b,
		})
	failOnError(err, "Fallo al publicar un mensaje")

	for d := range msgs {
		if corrId == d.CorrelationId {
			trans = bytetotrans(d.Body)
			break
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	nombre := nombreFrom(os.Args)

	fmt.Printf("Hola, mi nombre es : %s\n", nombre)
	ops := randInt(1, maxOp)
	fmt.Printf("%s quiere hacer %d operaciones\n", nombre, ops)
	for i := 1; i <= ops; i++ {
		cantidad := randInt(-10, +10)
		// La cantidad no puede ser 0
		for cantidad == 0 {
			cantidad = randInt(-10, +10)
		}
		fmt.Printf("%s operación %d: %d\n", nombre, i, cantidad)
		trans, err := operacionRPC(int32(cantidad))
		failOnError(err, "Fallo al solicitar la operación")
		fmt.Printf("Operación solicitada\n")
		if trans.permiso {
			fmt.Printf("Operación completada!\n")
		} else {
			fmt.Printf("NO PERMITIDA, NO HAY FONDOS\n")
		}
		fmt.Printf("Balance actual: %d\n", trans.balance)
		fmt.Printf("%d", i)
		fmt.Println(strings.Repeat("-", 30))
	}
}

func nombreFrom(args []string) string {
	if (len(args) < 2) || os.Args[1] == "" {
		return "Cliente"
	}
	return args[1]
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
