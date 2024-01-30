package cli

import (
	"encoding/json"
	"log"

	"go.uber.org/zap"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/pkg/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spf13/cobra"
)

type DeadLetterQ struct {
	Handler    string `json:"handler_poisoned"`
	Reason     string `json:"reason_poisoned"`
	Subscriber string `json:"subscriber_poisoned"`
	Topic      string `json:"topic_poisoned"`
}

// GetAPICommandDef runs app
func GetDeadQueueCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {

	workerCommand := cobra.Command{
		Use:   "dead-letter-queue",
		Short: "To start dead-letter queue",
		Long:  `This queue is used to store failed job from all worker`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Init worker
			subscriber, err := watermill.InitSubscriber(cfg, true)
			if err != nil {
				return err
			}

			// run worker with topic(queue name) and process function
			err = subscriber.Run(cfg.MQ.DeadLetterQ,"dead_letter_queue" ,HandleFailJob)
			return err
		},
	}
	return workerCommand
}

func HandleFailJob(msg *message.Message) error {
	// get fail job details from metadata
	var result DeadLetterQ
	metadata := make(map[string]string)
	for k, v := range msg.Metadata {
		metadata[k] = v
	}
	jsonString, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		return err
	}

	log.Printf("Failed job handler: %s, reason: %s, subscriber: %s, topic: %s", result.Handler, result.Reason, result.Subscriber, result.Topic)
	return nil
}
