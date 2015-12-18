package transport

import (
	"github.com/satori/go.uuid"
)

var ID = struct {
	Default int
	UUID    int
}{
	0, 1,
}

/*
 An object that can generate a series of identiy, typed as string.
 Each idenity should be unique.
*/
type IDMaker interface {
	// routine(thread) safe is required.
	NewID() string
}

// default IDMaker
type uuidMaker struct{}

func (*uuidMaker) NewID() string {
	return uuid.NewV4().String()
}
