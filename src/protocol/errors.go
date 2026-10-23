package protocol

const (
	// General Errors
	ErrCast = "cast error"

	// Validation Errors
	ErrInvalidName = "invalid name, must respect format " + ValidRegex
	ErrInvalidKey  = "invalid key"

	// Topic Errors
	ErrTopicAlreadyExist = "topic already exists"
	ErrTopicDoesNotExist = "topic does not exist"

	// Tree and Client Errors
	ErrSubscriberAlreadyActive = "subscriber is already active"
	ErrClientAlreadyInserted   = "client is already inserted inside the tree"
	ErrClientNotInserted       = "client is not inserted inside the tree"
)
