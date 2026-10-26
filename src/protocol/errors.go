package protocol

const (
	// General Errors
	ErrCast           = "cast error"
	ErrInvalidCommand = "invalid command. Need to be: create, list, produce, subscribe or info"

	// Validation Errors
	ErrInvalidName   = "invalid name, must respect format " + ValidRegex
	ErrInvalidFormat = "invalid format, must respect format 'key:body'"
	ErrInvalidKey    = "invalid key"

	// Topic Errors
	ErrTopicAlreadyExist = "topic already exists"
	ErrTopicDoesNotExist = "topic does not exist"

	// Tree and Client Errors
	ErrSubscriberAlreadyActive = "subscriber is already active"
	ErrClientAlreadyInserted   = "client is already inserted inside the tree"
	ErrClientNotInserted       = "client is not inserted inside the tree"
)

type Error struct {
	Message string
	Code    int
}
