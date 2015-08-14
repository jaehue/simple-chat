package common

import (
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/objx"
)

// State represents a map of state arguments that can be used to
// persist values across the authentication process.
type State struct {
	objx.Map
}

// NewState creates a new object that can be used to persist
// state across authentication requests.
func NewState(keyAndValuePairs ...interface{}) *State {
	return &State{objx.MSI(keyAndValuePairs...)}
}
