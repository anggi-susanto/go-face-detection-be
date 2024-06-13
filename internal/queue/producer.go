package queue

import (
	"fmt"

	"github.com/anggi-susanto/go-face-detection-be/config"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Producer interface {
	SendToQueue(photoID int64) error
}

type producer struct {
	config *config.RabbitMqConfig
}

func NewProducer(config *config.RabbitMqConfig) Producer {
	return &producer{
		config: config,
	}
}

func (p *producer) SendToQueue(photoID int64) error {
	conn, err := amqp.Dial(p.config.Uri)
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"face_detection",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatal(err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%d", photoID)),
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}
