package watermill

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/Improwised/golang-api/cli/workers"
	"github.com/Improwised/golang-api/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"

	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type WatermillPubliser struct {
	publisher message.Publisher
}

func InitPubliser(cfg config.AppConfig) (*WatermillPubliser, error) {
	logger = watermill.NewStdLogger(cfg.MQ.Debug, cfg.MQ.Track)
	switch cfg.MQ.Dialect {
	case "amqp":
		return initAmqpPub(cfg)

	case "redis":
		return initRedisPub(cfg)

	case "kafka":
		return initKafkaPub(cfg)

	default:
		return &WatermillPubliser{}, nil
	}

}

// send message to queue using topic name
//
// struct must from worker package(/cli/workers)
func (wp *WatermillPubliser) Publish(topic string, data interface{}) error {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	var handle workers.Handler

	handle, ok := data.(workers.Handler)
	if !ok {
		return fmt.Errorf("data is not of type workers.Handler")
	}

	err := enc.Encode(&handle)
	if err != nil {
		log.Fatal("encode: niche", err)
	}
	msg := message.NewMessage(watermill.NewUUID(), network.Bytes())

	if err := wp.publisher.Publish(topic, msg); err != nil {
		return err
	}
	return nil
}

func initAmqpPub(cfg config.AppConfig) (*WatermillPubliser, error) {
	amqpConfig := amqp.NewDurableQueueConfig(cfg.MQ.Amqp.AmqbUrl)
	publisher, err := amqp.NewPublisher(amqpConfig, logger)
	return &WatermillPubliser{publisher: publisher}, err
}

// TODO: username/pass
func initRedisPub(cfg config.AppConfig) (*WatermillPubliser, error) {
	pubClient := redis.NewClient(&redis.Options{
		Addr: cfg.MQ.Redis.RedisUrl,
	})
	publisher, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client:     pubClient,
			Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		logger,
	)
	return &WatermillPubliser{publisher: publisher}, err
}

func initKafkaPub(cfg config.AppConfig) (*WatermillPubliser, error) {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:               cfg.MQ.Kafka.KafkaBroker,
			Marshaler:             kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: kafka.DefaultSaramaSyncPublisherConfig(),
		},
		logger,
	)
	return &WatermillPubliser{publisher: publisher}, err
}

