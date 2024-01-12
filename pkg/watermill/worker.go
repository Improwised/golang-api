package watermill

import (
	"context"
	"log"
	"time"

	"github.com/Improwised/golang-api/config"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"

	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/redis/go-redis/v9"
)

var logger = watermill.NewStdLogger(true, false)

type WatermillSubscriber struct {
	Subscriber message.Subscriber
	Router     *message.Router
}

func InitWorker(cfg config.AppConfig) (*WatermillSubscriber, error) {
	switch cfg.MQDialect {
	case "amqp":
		amqpConfig := amqp.NewDurableQueueConfig(cfg.AMQB.AmqbUrl)
		subscriber, err := amqp.NewSubscriber(
			amqpConfig,
			watermill.NewStdLogger(false, false),
		)
		return &WatermillSubscriber{Subscriber: subscriber}, err

	case "redis":
		subClient := redis.NewClient(&redis.Options{
			Addr: cfg.Redis.RedisUrl,
		})
		subscriber, err := redisstream.NewSubscriber(
			redisstream.SubscriberConfig{
				Client:        subClient,
				Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
				ConsumerGroup: cfg.Redis.ConsumerGroup,
			},
			watermill.NewStdLogger(false, false),
		)
		return &WatermillSubscriber{Subscriber: subscriber}, err

	case "kafka":
		saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
		saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
		subscriber, err := kafka.NewSubscriber(
			kafka.SubscriberConfig{
				Brokers:               cfg.Kafka.KafkaBroker,
				Unmarshaler:           kafka.DefaultMarshaler{},
				OverwriteSaramaConfig: saramaSubscriberConfig,
				ConsumerGroup:         cfg.Kafka.ConsumerGroup,
			},
			watermill.NewStdLogger(false, false),
		)
		if err != nil {
			return nil, err
		}
		return &WatermillSubscriber{Subscriber: subscriber}, err
	default:
		return nil, nil
	}
}

func (ws *WatermillSubscriber) InitRouter(cfg config.AppConfig, delayTime, MaxRetry int) (*WatermillSubscriber, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	pub, err := InitSender(cfg)
	if err != nil {
		return nil, err
	}

	poq, err := middleware.PoisonQueue(pub.publisher, "poison_queue")
	if err != nil {
		return nil, err
	}
	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		poq,
		middleware.Retry{
			MaxRetries:      MaxRetry,
			Logger:          logger,
			MaxInterval:     time.Millisecond * time.Duration(delayTime),
			InitialInterval: time.Millisecond * time.Duration(delayTime),
			Multiplier:      1,
			OnRetryHook: func(retryNum int, delay time.Duration) {
				log.Println("retry count :=", delay)
			},
		}.Middleware,

		middleware.Recoverer,
	)
	ws.Router = router

	return ws, nil
}

func (ws *WatermillSubscriber) Run(topic string, handlerFunc message.NoPublishHandlerFunc) error {
	ws.Router.AddNoPublisherHandler(
		"counter",
		topic,
		ws.Subscriber,
		handlerFunc,
	)

	err := ws.Router.Run(context.Background())
	return err
}
