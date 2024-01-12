package watermill

import (
	"github.com/Improwised/golang-api/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"

	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

// GetAPICommandDef runs app
type WatermillPubliser struct {
	publisher message.Publisher
}

func InitSender(cfg config.AppConfig) (*WatermillPubliser, error) {
	switch cfg.MQDialect {
	case "amqp":
		amqpConfig := amqp.NewDurableQueueConfig(cfg.AMQB.AmqbUrl)
		publisher, err := amqp.NewPublisher(amqpConfig, watermill.NewStdLogger(false, false))
		return &WatermillPubliser{publisher: publisher}, err

	case "redis":
		pubClient := redis.NewClient(&redis.Options{
			Addr: cfg.Redis.RedisUrl,
		})
		publisher, err := redisstream.NewPublisher(
			redisstream.PublisherConfig{
				Client:     pubClient,
				Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
			},
			watermill.NewStdLogger(false, false),
		)

		return &WatermillPubliser{publisher: publisher}, err

	case "kafka":
		publisher, err := kafka.NewPublisher(
			kafka.PublisherConfig{
				Brokers:               cfg.Kafka.KafkaBroker,
				Marshaler:             kafka.DefaultMarshaler{},
				OverwriteSaramaConfig: kafka.DefaultSaramaSyncPublisherConfig(),
			},
			watermill.NewStdLogger(false, false),
		)
		return &WatermillPubliser{publisher: publisher}, err

	default:
		return &WatermillPubliser{}, nil
	}

}
func (wp *WatermillPubliser) PublishMessages(topic string, data []byte) error {
	msg := message.NewMessage(watermill.NewUUID(), data)
	if err := wp.publisher.Publish(topic, msg); err != nil {
		return err
	}
	return nil
}
