package config

type Config struct {
	MongoConfig    MongoConfig
	RabbitMqConfig RabbitMqConfig
}

type MongoConfig struct {
	Uri        string
	Database   string
	Collection string
}

type RabbitMqConfig struct {
	Uri string
}
