package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	color := randNombre()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fallo al conectar a RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Fallo al abrir un canal")
	defer ch.Close()

	// Cola para ser avisado por el banquero
	ql, err := ch.QueueDeclare(
		"avisos", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Fallo al declarar una cola")

	msgs, err := ch.Consume(
		ql.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Fallo al registrar un consumidor")
	log.Printf(" [*] Esperando los avisos. Para salir pulsa CTRL+C\n")
	forever := make(chan bool)
	fmt.Printf("Hola, mi nombre es: Señor %s\n", color)
	go func() {
		for d := range msgs {
			fmt.Printf("¡Esto es un atraco, por menos de 20 no me muevo de aquí!\n")
			time.Sleep(2 * time.Second)
			fmt.Printf("El botín es de %d\n", int32(binary.LittleEndian.Uint32(d.Body)))
			fmt.Printf("Me voy corriendo\n")
		}
	}()
	<-forever
}

func randNombre() string {
	colores := [6]string{"Marrón", "Rubio", "Rosa", "Blanco", "Naranja", "Azul"}
	return colores[rand.Intn(len(colores)-1)]
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
