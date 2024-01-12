package config

type RedisConfig struct {
	RedisUrl string `envconfig:"REDIS_URI"`
	ConsumerGroup string `envconfig:"CONSUMER_GROUP"`
}

type AmqpConfig struct {
	AmqbUrl string `envconfig:"AMQB_URI"`
}
type KafkaConfig struct {
	KafkaBroker []string `envconfig:"KAFKA_BROKER"`
	ConsumerGroup string `envconfig:"CONSUMER_GROUP"`
}
