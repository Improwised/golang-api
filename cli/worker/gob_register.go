package worker

import (
	"encoding/gob"
)

func GobRegister() {
	gob.Register(UpdateUser{})
    gob.Register(DeleteUser{})
    gob.Register(AddUser{})
}
