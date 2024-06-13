package queue

import (
	"context"
	"fmt"
	"strconv"

	"github.com/anggi-susanto/go-face-detection-be/config"
	"github.com/anggi-susanto/go-face-detection-be/domain"
	"github.com/anggi-susanto/go-face-detection-be/internal/repository/mongo"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	python3 "github.com/DataDog/go-python3"
)

type Consumer interface {
	ReceiveFromQueue(ctx context.Context)
}

type consumer struct {
	config *config.RabbitMqConfig
	repo   *mongo.PhotoRepository
}

func NewConsumer(config *config.RabbitMqConfig, repo *mongo.PhotoRepository) Consumer {
	return &consumer{
		config: config,
		repo:   repo,
	}
}

func (c *consumer) ReceiveFromQueue(ctx context.Context) {
	conn, err := amqp.Dial(c.config.Uri)
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

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			photoIDStr := string(d.Body)

			photo, err := c.repo.FindByID(ctx, photoIDStr)
			if err != nil {
				logrus.Infof("Error retrieving file path: %s", err)
				c.updateStatus(ctx, photo, "error", 0)
				continue
			}

			facesDetected, err := DetectFaces(photo.FilePath)
			if err != nil {
				logrus.Infof("Error processing photo: %s", err)
				c.updateStatus(ctx, photo, "error", 0)
			} else {
				logrus.Infof("Successfully processed photo: %s. Faces detected: %d", photo.FilePath, facesDetected)
				c.updateStatus(ctx, photo, "processed", facesDetected)
			}
		}
	}()

	logrus.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c *consumer) updateStatus(ctx context.Context, photo *domain.Photo, status string, facesDetected int) {
	photo.FacesDetected = facesDetected
	photo.Status = status
	err := c.repo.Update(ctx, photo)
	if err != nil {
		logrus.Errorf("Failed to update status: %v", err)
	}
}

func DetectFaces(imagePath string) (int, error) {
	python3.Py_Initialize()
	defer python3.Py_Finalize()

	code := `
import sys
import cv2

def detect_faces(image_path):
    img = cv2.imread(image_path)
    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    faces = face_cascade.detectMultiScale(gray, 1.1, 4)
    
    for (x, y, w, h) in faces:
        cv2.rectangle(img, (x, y), (x+w, y+h), (255, 0, 0), 2)
    
    result_path = image_path.replace("/app/photos/", "/app/photos/processed_")
    cv2.imwrite(result_path, img)
    return len(faces)
`

	pythonModule := python3.PyImport_AddModule("__main__")

	python3.PyRun_SimpleString(code)

	detectFaces := pythonModule.GetAttrString("detect_faces")
	if detectFaces == nil {
		return 0, fmt.Errorf("Failed to load function detect_faces")
	}

	args := python3.PyTuple_New(1)
	python3.PyTuple_SetItem(args, 0, python3.PyUnicode_FromString(imagePath))

	result := detectFaces.CallObject(args)
	if result == nil {
		return 0, fmt.Errorf("Failed to call function detect_faces")
	}

	facesDetected := python3.PyLong_AsLong(result)
	return facesDetected, nil
}

func (c *consumer) getFilePath(ctx context.Context, photoID int64) (string, error) {
	photo, err := c.repo.FindByID(ctx, strconv.FormatInt(photoID, 10))
	if err != nil {
		return "", err
	}

	return photo.FilePath, nil
}
