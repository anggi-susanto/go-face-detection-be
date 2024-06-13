package main

import (
	"context"
	"os"
	"time"

	"github.com/anggi-susanto/go-face-detection-be/config"
	"github.com/anggi-susanto/go-face-detection-be/internal/queue"
	mongoRepo "github.com/anggi-susanto/go-face-detection-be/internal/repository/mongo"
	"github.com/anggi-susanto/go-face-detection-be/internal/rest"
	"github.com/anggi-susanto/go-face-detection-be/photo"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func initMongo(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		logrus.Fatalf("Failed to ping MongoDB: %v", err)
	}

	MongoClient = client
	logrus.Println("Connected to MongoDB")
}

func main() {

	// setup config
	config := config.Config{
		MongoConfig: config.MongoConfig{
			Uri:        os.Getenv("MONGO_URI"),
			Database:   os.Getenv("MONGO_DB"),
			Collection: os.Getenv("MONGO_COLLECTION"),
		},
		RabbitMqConfig: config.RabbitMqConfig{
			Uri: os.Getenv("RABBITMQ_URI"),
		},
	}

	initMongo(config.MongoConfig.Uri)
	defer MongoClient.Disconnect(context.Background())

	// setup fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Get("/docs/*", swagger.HandlerDefault)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("MRT API is UP and RUNNING!")
	})

	// handler composing
	photoRepo := mongoRepo.NewPhotoRepository(MongoClient, &config.MongoConfig)
	photoService := photo.NewService(*photoRepo)
	photoProducer := queue.NewProducer(&config.RabbitMqConfig)
	photoHandler := rest.NewPhotoHandler(photoService, photoProducer)

	// route definitions
	app.Post("/upload", photoHandler.Upload)
	app.Get("/result/:id", photoHandler.CheckResult)
	app.Get("/photo/:id", photoHandler.GetPhoto)

	// consumer starting up
	consumer := queue.NewConsumer(&config.RabbitMqConfig, photoRepo)
	consumer.ReceiveFromQueue(context.Background())

	app.Listen(":8080")
}
