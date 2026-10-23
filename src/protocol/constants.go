package protocol

const (
	EntryPoint = "/tmp/operator.server."

	FnvPrime       = 1099511628211
	FnvOffsetBasis = 14695981039346656037

	MaxIdCharacters    = 32
	ClientsHashmapSize = 64

	ValidRegex = `^[a-zA-Z0-9_.-]{1,32}$`
)
