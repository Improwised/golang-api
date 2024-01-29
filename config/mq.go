package config

type MQConfig struct {
	Dialect     string `envconfig:"MQ_DIALECT"`
	Debug       bool   `envconfig:"MQ_DEBUG"`
	Track       bool   `envconfig:"MQ_TRACK"`
	DeadQueue   string `envconfig:"DEAD_LETTER_QUEUE"`
	Redis       RedisConfig
	Amqp        AmqpConfig
	Kafka       KafkaConfig
	GoogleCloud GoogleCloud
}
type RedisConfig struct {
	RedisUrl      string `envconfig:"REDIS_URI"`
	ConsumerGroup string `envconfig:"CONSUMER_GROUP"`
	UserName      string `envconfig:"REDIS_USERNAME"`
	Password      string `envconfig:"REDIS_PASSWORD"`
}

type AmqpConfig struct {
	AmqbUrl string `envconfig:"AMQB_URI"`
}
type KafkaConfig struct {
	KafkaBroker   []string `envconfig:"KAFKA_BROKER"`
	ConsumerGroup string   `envconfig:"CONSUMER_GROUP"`
}

type GoogleCloud struct {
	ProjectID      string `envconfig:"PROJECT_ID"`
	SubscriptionId string `envconfig:"SUBSCRIPTION_ID"`
}
