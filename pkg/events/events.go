package events

import (
	"log"

	"github.com/Improwised/golang-api/constants"
	"github.com/Improwised/golang-api/pkg/structs"
	evbus "github.com/asaskevich/EventBus"
)

//go:generate mockery --name IEvents
type IEvents interface {
	SubscribeUserRegistered() error
	Publish(event string, data interface{})
}

type Events struct {
	Bus evbus.Bus
}

func NewEventBus() *Events {
	return &Events{
		Bus: evbus.New(),
	}
}

func userRegistered(data structs.EventUserRegistered) {
	log.Printf("User Registered email: %s", data.Email)
	// We can send message to SQS, Redis, or any other place
}

func (eb *Events) SubscribeUserRegistered() error {
	return eb.Bus.Subscribe(constants.EventUserRegistered, userRegistered)
}

func (eb *Events) Publish(event string, data interface{}) {
	eb.Bus.Publish(event, data)
}

// SubscribeAll start all subscribers
func (eb *Events) SubscribeAll() error {
	err := eb.Bus.Subscribe(constants.EventUserRegistered, userRegistered)
	if err != nil {
		return err
	}
	return nil
}
