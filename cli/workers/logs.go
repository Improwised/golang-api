package workers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"
	"github.com/Improwised/golang-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type JsonLogs struct {
	Level      string `json:"level"`
	Timestamp  string `json:"timestamp"`
	Caller     string `json:"caller"`
	Message    string `json:"message"`
	Host       string `json:"host"`
	Method     string `json:"method"`
	UserId     string `json:"user_id"`
	StatusCode int    `json:"status"`
}

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

func (j JsonLogs) Handle() error {
	jsonData, err := json.Marshal(j)
	if err != nil {
		return err
	}

	log.Printf("%s\n", jsonData)
	return nil
}

func (e EventLogs) Handle() error {
	log.Printf("inserting event logs in database")
	cfg := config.GetConfig()
	db, err := database.Connect(cfg.DB)
	if err != nil {
		return err
	}
	
	logs, err := models.InitLogModel(db)
	if err != nil {
		return err
	}
	
	eventLog := models.EventLogs(e)
	err=logs.CreateLog(eventLog)
	if err != nil {
		return err
	}
	log.Printf("logs inserted successfully")
	return nil
}

func GetJsonLogs(c *fiber.Ctx, caller, level, userID, message string, statusCode int) JsonLogs {
	return JsonLogs{
		Level:      level,
		Method:     c.Method(),
		Timestamp:  time.Now().Format(time.RFC3339),
		Caller:     caller,
		Host:       c.Hostname(),
		StatusCode: statusCode,
		UserId:     userID,
		Message:    message,
	}
}

func GetEventLogs(c *fiber.Ctx, caller, level, userID, message string, statusCode int) EventLogs {
	return EventLogs{
		UUID:       uuid.New().String(),
		Level:      level,
		Method:     c.Method(),
		Caller:     caller,
		Host:       c.Hostname(),
		StatusCode: statusCode,
		UserId:     userID,
		Message:    message,
		Payload:    string(c.Body()),
	}
}
