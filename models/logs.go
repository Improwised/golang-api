package models

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

type EventLogs struct {
	UUID       string `json:"uuid"`
	Level      string `json:"level"`
	Caller     string `json:"caller"`
	Message    string `json:"message"`
	Host       string `json:"host"`
	Method     string `json:"method"`
	UserId     string `json:"user_id"`
	StatusCode int    `json:"status"`
	Payload    string `json:"payload"`
}

type LogModel struct {
	db *goqu.Database
}

const logTable = "event_logs"

// InitUserModel Init model
func InitLogModel(goqu *goqu.Database) (LogModel, error) {
	return LogModel{
		db: goqu,
	}, nil
}

func (model *LogModel) CreateLog(log EventLogs) error {
	_, err := model.db.Insert(logTable).Rows(
		goqu.Record{
			"uuid":       log.UUID,
			"level":      log.Level,
			"created_at": time.Now().Format(time.RFC3339),
			"caller":     log.Caller,
			"message":    log.Message,
			"host":       log.Host,
			"method":     log.Method,
			"user_id":    log.UserId,
			"status":     log.StatusCode,
			"payload":    log.Payload,
		},
	).Executor().Exec()
	if err != nil {
		return err
	}
	return nil
}
