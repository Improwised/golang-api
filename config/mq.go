package config

// TODO: add env to .env.example
type MQConfig struct {
	Dialect   string `envconfig:"MQ_DIALECT"`
	Debug     bool   `envconfig:"MQ_DEBUG"`
	Track     bool   `envconfig:"MQ_TRACK"`
	DeadQueue string `envconfig:"DEAD_LETTER_QUEUE"`
	Redis     RedisConfig
	Amqp      AmqpConfig
	Kafka     KafkaConfig
}
type RedisConfig struct {
	RedisUrl      string `envconfig:"REDIS_URI"`
	ConsumerGroup string `envconfig:"CONSUMER_GROUP"`
}

type AmqpConfig struct {
	AmqbUrl string `envconfig:"AMQB_URI"`
}
type KafkaConfig struct {
	KafkaBroker   []string `envconfig:"KAFKA_BROKER"`
	ConsumerGroup string   `envconfig:"CONSUMER_GROUP"`
}
