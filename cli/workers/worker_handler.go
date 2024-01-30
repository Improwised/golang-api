package workers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
)

func init() {
	for _, v := range RegisterWorkerStruct() {
		gob.Register(v)
	}
}

// Register all worker struct here befour run worker for proper unmarshalling
func RegisterWorkerStruct() []interface{} {
	return []interface{}{
		WelcomeMail{},
		// ...
	}
}

// Handler interface for all worker struct
type Handler interface {
	Handle() error
}

// process all worker struct and call Handle function according to struct
func Process(msg *message.Message) error {
	buf := bytes.NewBuffer(msg.Payload)
	dec := gob.NewDecoder(buf)

	var result Handler
	err := dec.Decode(&result)
	if err != nil {
		return err
	}
	// Store the JSON payload in the msg so that it can be persisted in the database in case the job fails.
	msg.Payload, err = json.Marshal(result)
	if err != nil {
		return err
	}
	if err := result.Handle(); err != nil {
		return err
	}
	msg.Ack()
	return nil
}
