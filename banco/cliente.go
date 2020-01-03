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

// Número máximo de operaciones bancarias
const maxOp = 8

// Tipo que se envía como resultado de la operación
type transaccion struct {
	permiso bool  // Resultado de la operacion
	balance int32 // Balance de la cuenta
}

// Pasa la secuencia de bytes a una transacción
func bytetotrans(arr []byte) transaccion {
	// Byte a int8 a boolean
	permiso := int8(arr[0]) != 0
	// 4 bytes a uint32 a int32
	balance := int32(binary.LittleEndian.Uint32(arr[1:]))
	return transaccion{permiso, balance}
}

// operacionRPC envía una operación bancaria al banquero
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

	// Se declara consumidor de la cola de respuestas del banquero
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

	// Correlation ID aleatorio que se usa para recibir respuesta
	corrId := randomString(32)
	// Mensaje que contiene la cantidad a ingresar o retirar
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(cantidad))
	// Se envía el mensaje
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

	// Se espera la respuesta del banquero
	for d := range msgs {
		if corrId == d.CorrelationId {
			trans = bytetotrans(d.Body)
			break
		}
	}
	// Se devuelve la transaccion y si ha habido errores
	return
}

func main() {
	// Seed para la librería math/rand con la fecha y hora actual
	rand.Seed(time.Now().UTC().UnixNano())

	// Se coge el nombre indicado o uno por defecto
	nombre := nombreFrom(os.Args)

	fmt.Printf("Hola, mi nombre es : %s\n", nombre)
	// Número aleatorio de operaciones entre 1 y el máximo de operaciones
	ops := randInt(1, maxOp)
	fmt.Printf("%s quiere hacer %d operaciones\n", nombre, ops)
	for i := 1; i <= ops; i++ {
		// Cantidad aleatoria a ingresar o retirar
		cantidad := randInt(-10, +10)
		// La cantidad no puede ser 0
		for cantidad == 0 {
			// Mientras la cantidad sea 0 se vuelve a escoger otra cantidad
			cantidad = randInt(-10, +10)
		}
		fmt.Printf("%s operación %d: %d\n", nombre, i, cantidad)
		// Se envía la operación al banquero
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
		// Espera antes de realizar otra operación
		time.Sleep(1 * time.Second)
	}
}

func nombreFrom(args []string) string {
	// Si no se ha indicado nombre, se le asigna "Cliente"
	if (len(args) < 2) || os.Args[1] == "" {
		return "Cliente"
	}
	// Devuelve el nombre asignado
	return args[1]
}

// randomString devuelve un String aleatorio de longitud "l"
func randomString(l int) string {
	// Array de bytes de tamaño "l"
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		// Asigna a cada byte el valor de una letra mayúscula aleatoria
		bytes[i] = byte(randInt(65, 90))
	}
	// Devuelve el string
	return string(bytes)
}

// randInt devuelve un entero aleatorio entre "min" y "max"
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// failOnError imprime un mensaje de error si ha habido algún error
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
