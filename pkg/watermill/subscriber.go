package watermill

import (
	"context"
	"fmt"
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

var logger watermill.LoggerAdapter

type WatermillSubscriber struct {
	Subscriber message.Subscriber
	Router     *message.Router
}


func InitSubscriber(cfg config.AppConfig) (*WatermillSubscriber, error) {
	logger = watermill.NewStdLogger(cfg.MQ.Debug, cfg.MQ.Track)
	switch cfg.MQ.Dialect {
	case "amqp":
		return initAmqpSub(cfg)

	case "redis":
		return initRedisSub(cfg)

	case "kafka":
		return initKafkaSub(cfg)

	default:
		return nil, nil
	}
}

// InitRouter init router for add middleware,retry count,delay etc
func (ws *WatermillSubscriber) InitRouter(cfg config.AppConfig, delayTime, MaxRetry int) (*WatermillSubscriber, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	pub, err := InitPubliser(cfg)
	if err != nil {
		return nil, err
	}

	poq, err := middleware.PoisonQueue(pub.publisher, cfg.MQ.DeadQueue)
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
		}.Middleware,

		middleware.Recoverer,
	)
	ws.Router = router

	return ws, nil
}

func (ws *WatermillSubscriber) Run(topic string, handlerFunc message.NoPublishHandlerFunc) error {
	if ws.Subscriber == nil {
		return fmt.Errorf("subscriber is nil")
	}

	if ws.Router == nil {
		router, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			return err
		}
		ws.Router = router
	}

	ws.Router.AddNoPublisherHandler(
		"handler",
		topic,
		ws.Subscriber,
		handlerFunc,
	)

	err := ws.Router.Run(context.Background())
	return err
}

func initAmqpSub(cfg config.AppConfig) (*WatermillSubscriber, error) {
	amqpConfig := amqp.NewDurableQueueConfig(cfg.MQ.Amqp.AmqbUrl)
	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		logger,
	)
	return &WatermillSubscriber{Subscriber: subscriber}, err
}

func initKafkaSub(cfg config.AppConfig) (*WatermillSubscriber, error) {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               cfg.MQ.Kafka.KafkaBroker,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         cfg.MQ.Kafka.ConsumerGroup,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}
	return &WatermillSubscriber{Subscriber: subscriber}, err
}

func initRedisSub(cfg config.AppConfig) (*WatermillSubscriber, error) {
	subClient := redis.NewClient(&redis.Options{
		Addr: cfg.MQ.Redis.RedisUrl,
		Username: cfg.MQ.Redis.UserName,
		Password: cfg.MQ.Redis.Password,
	})
	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        subClient,
			Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
			ConsumerGroup: cfg.MQ.Redis.ConsumerGroup,
		},
		logger,
	)
	return &WatermillSubscriber{Subscriber: subscriber}, err
}
