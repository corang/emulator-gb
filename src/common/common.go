package common

import (
	"errors"
)

func Unimplemented() {
	// Unimplemented
	panic(errors.New("Unimplemented"))
}
