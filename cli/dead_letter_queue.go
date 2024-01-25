package cli

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/pkg/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spf13/cobra"
)

// GetAPICommandDef runs app
func GetDeadQueueCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {

	workerCommand := cobra.Command{
		Use:   "dead-letter-queue",
		Short: "To start dead-letter queue",
		Long:  `This queue is used to store failed job from all worker`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Init worker
			subscriber, err := watermill.InitSubscriber(cfg)
			if err != nil {
				return err
			}

			// run worker with topic(queue name) and process function
			err = subscriber.Run(cfg.MQ.DeadQueue, HandleFailJob)
			return err
		},
	}
	return workerCommand
}

func HandleFailJob(msg *message.Message) error {
	fmt.Println("failed job:-", msg.UUID)
	// process here
	return nil
}
