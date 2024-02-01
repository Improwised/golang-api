package workers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strconv"

	"github.com/ThreeDotsLabs/watermill/message"
)

type QConfig struct {
	Topic      string
	MaxRetries int
}

func init() {
	for _, v := range RegisterWorkerStruct() {
		gob.Register(v)
	}
}

// Register all worker struct here befour run worker for proper unmarshalling
func RegisterWorkerStruct() []interface{} {
	return []interface{}{
		WelcomeMail{},
		LoginMail{},
		JsonLogs{},
		EventLogs{},
		// ...
	}
}

// Handler interface for all worker struct
type Handler interface {
	Handle() error
}

// process all worker struct and call Handle function according to struct
func (q QConfig) Process(msg *message.Message) error {
	buf := bytes.NewBuffer(msg.Payload)
	dec := gob.NewDecoder(buf)

	var result Handler
	err := dec.Decode(&result)
	if err != nil {
		return err
	}

	// Store the JSON payload in the msg so that it can be persisted in the database in case the job fails.
	CurrentRetryCount := GetSetRetryCount(msg)
	if CurrentRetryCount == q.MaxRetries {
		msg.Payload, err = json.Marshal(result)
		if err != nil {
			return err
		}
	}

	if err := result.Handle(); err != nil {
		return err
	}
	msg.Ack()
	return nil
}

// set/get retry count from message metadata
func GetSetRetryCount(msg *message.Message) int {
	count := msg.Metadata.Get("max_retry_count")
	if count == "" {
		msg.Metadata.Set("max_retry_count", "0")
		count = "0"
		return 0
	}
	countInt, _ := strconv.Atoi(count)
	msg.Metadata.Set("max_retry_count", strconv.Itoa(countInt+1))
	return countInt + 1
}
