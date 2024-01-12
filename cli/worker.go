package cli

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"
	"github.com/Improwised/golang-api/models"
	"github.com/Improwised/golang-api/pkg/events"
	"github.com/Improwised/golang-api/pkg/structs"
	"github.com/Improwised/golang-api/pkg/watermill"
	"github.com/Improwised/golang-api/services"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spf13/cobra"
)

// GetAPICommandDef runs app

type worker struct {
	userModel *services.UserService
	events    *events.Events
}

func GetWorkerCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {

	workerCommand := cobra.Command{
		Use:   "worker",
		Short: "To start worker",
		Long:  `To start worker`,
		RunE: func(cmd *cobra.Command, args []string) error {

			events := events.NewEventBus(logger)

			goqu, err := database.Connect(cfg.DB)
			if err != nil {
				return err
			}
			err = events.SubscribeAll()
			if err != nil {
				return err
			}

			userModel, err := models.InitUserModel(goqu)
			if err != nil {
				return err
			}
			userSvc := services.NewUserService(&userModel)

			// Init worker
			sub, err := watermill.InitWorker(cfg)
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
			router, err := sub.InitRouter(cfg, retryCount, delay)
			if err != nil {
				return err
			}

			w := worker{
				userModel: userSvc,
				events:    events,
			}
			switch topic {
			case "create_user":
				err = router.Run(topic, w.Process)
				return err
			}
			return nil
		},
	}

	return workerCommand
}

func (w *worker) Process(msg *message.Message) error {
	var userReq structs.ReqRegisterUser

	err := json.Unmarshal(msg.Payload, &userReq)
	if err != nil {
		return err
	}
	_, err = w.userModel.RegisterUser(models.User{FirstName: userReq.FirstName, LastName: userReq.LastName, Email: userReq.Email, Password: userReq.Password, Roles: userReq.Roles}, w.events)

	return err
}
