package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

// Balance de la cuenta compartida
var balance int32

// Unidades necesarias para avisar al ladrón
const avisoLadron = 20

// Tipo que se envía como resultado de la operación
type transaccion struct {
	permiso bool  // Resultado de la operacion
	balance int32 // Balance de la cuenta
}

// transtobyte pasa una transacción a una secuencia de bytes
func transtobyte(trans transaccion) []byte {
	// Array de 5 bytes
	b := make([]byte, 5)
	// Pasa permiso de boolean a int8 a byte
	b[0] = byte(Btoi(trans.permiso))
	// Pasa el balance de int32 a uint32 a secuencia de 4 bytes
	binary.LittleEndian.PutUint32(b[1:], uint32(trans.balance))
	return b
}

// Btoi pasa un bool a int
func Btoi(b bool) int8 {
	if b {
		// Si es true devuelve 1
		return 1
	}
	// Si es false devuelve 0
	return 0
}

// operacion intenta realizar una operación bancaria
func operacion(n int32) bool {
	// Simula el tiempo que tarda en hacer la operación
	time.Sleep(1 * time.Second)
	// Si no hay fondos suficientes devuelve false
	if (balance + n) < 0 {
		return false
	}
	// Si hay fondos suficientes devuelve true y actualiza el balance
	balance += n
	return true
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fallo al conectar a RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Fallo al abrir un canal")
	defer ch.Close()

	// Cola para las operaciones bancarias
	q, err := ch.QueueDeclare(
		"operaciones", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Fallo al declarar una cola")

	// Cola para avisar al ladrón
	ql, err := ch.QueueDeclare(
		"avisos", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Fallo al declarar una cola")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Fallo al configurar QoS")

	// Se declara consumidor de la cola de operaciones bancarias
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
		// Espera a que le lleguen operaciones bancarias
		for d := range msgs {
			fmt.Println(strings.Repeat(">", 30))
			// Pasa la cantidad de secuencia de bytes a uint32 a int32
			cantidad := int32(binary.LittleEndian.Uint32(d.Body))
			fmt.Printf("Operación recibida: %d\n", cantidad)
			permiso := operacion(cantidad)
			if !permiso {
				fmt.Printf("NO PERMITIDO, NO HAY FONDOS\n")
			}
			if balance >= avisoLadron {
				// Avisar al ladrón
				b := make([]byte, 4)
				// Pasa el botín de int 32 a uint32 a secuencia de bytes
				binary.LittleEndian.PutUint32(b, uint32(balance))

				// Envía el aviso con el botín al ladrón
				err = ch.Publish(
					"",      // exchange
					ql.Name, // routing key
					false,   // mandatory
					false,   // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(b),
					})
				failOnError(err, "Fallo al publicar un mensaje")
				// Después del robo el balance de la cuenta es 0
				balance = 0
			}
			// Indica el balance de la cuenta
			fmt.Printf("Balance: %d\n", balance)

			// Envía el resultado de la transacción al cliente
			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					// Pasa la transaccion a una secuencia de bytes
					Body: transtobyte(transaccion{permiso, balance}),
				})
			failOnError(err, "Fallo al publicar un mensaje")
			fmt.Println(strings.Repeat("<", 30))

			d.Ack(false)
		}
	}()

	// El banquero sigue en marcha hasta que se pulse CTRL+C
	<-forever
}

// failOnError imprime un mensaje de error si ha habido algún error
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
