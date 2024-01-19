package cli

import (
	"go.uber.org/zap"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/pkg/watermill"
	"github.com/Improwised/golang-api/cli/worker"
	"github.com/spf13/cobra"
)

// GetAPICommandDef runs app

func GetWorkerCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {

	workerCommand := cobra.Command{
		Use:   "worker",
		Short: "To start worker",
		Long:  `To start worker`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Init worker
			subscriber, err := watermill.InitWorker(cfg)
			if err != nil {
				return err
			}

			// get topic from flag
			topic, err := cmd.Flags().GetString("topic")
			if err != nil {
				return err
			}

			// get retry count from flag
			retryCount, err := cmd.Flags().GetInt("retry-count")
			if err != nil {
				return err
			}

			// get delay from flag
			delay, err := cmd.Flags().GetInt("delay")
			if err != nil {
				return err
			}

			// Init router for add middleware,retry count,etc
			router, err := subscriber.InitRouter(cfg, delay, retryCount)
			if err != nil {
				return err
			}

			err = router.Run(topic, worker.Process)
			return err

		},
	}

	return workerCommand
}
