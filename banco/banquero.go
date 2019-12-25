package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

// Balance de la cuenta compartida
var balance int32

// Unidades necesarias para avisar al ladrón
const avisoLadron = 20

type transaccion struct {
	permiso bool
	balance int32
}

// transtobyte pasa una transacción a una secuencia de bytes
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

func operacion(n int32) transaccion {
	if (balance + n) < 0 {
		return transaccion{false, balance}
	}
	balance = balance + n
	return transaccion{true, balance}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fallo al conectar a RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Fallo al abrir un canal")
	defer ch.Close()

	// Cola única del banco para las operaciones
	q, err := ch.QueueDeclare(
		"operaciones", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Fallo al declarar una cola")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Fallo al configurar QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Fallo al registrar un consumidor")

	forever := make(chan bool)
	fmt.Printf("El banco abre\n")
	log.Printf(" [*] Esperando clientes. Para salir pulsa CTRL+C\n")
	go func() {
		for d := range msgs {
			cantidad := int32(binary.LittleEndian.Uint32(d.Body))
			fmt.Println(strings.Repeat(">", 30))
			fmt.Printf("Operación recibida: %d\n", cantidad)
			respuesta := operacion(cantidad)
			if !respuesta.permiso {
				fmt.Printf("NO PERMITIDO, NO HAY FONDOS\n")
			}
			fmt.Printf("Balance: %d\n", respuesta.balance)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          transtobyte(respuesta),
				})
			failOnError(err, "Fallo al publicar un mensaje")
			fmt.Println(strings.Repeat("<", 30))

			d.Ack(false)
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
