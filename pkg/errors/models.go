package errors

// common errors
var (
	ErrIndexOutOfBounds    = Errorf("index out of bounds")
	ErrInvalidSliceLength  = Errorf("slice has invalid length")
	ErrInterfaceConversion = Errorf("failed to convert interface")
	ErrMissingContext      = Errorf("failed to find context")
	ErrNilInterface        = Errorf("interface is nil")
	ErrMissingMapKey       = Errorf("key not found in map")
)
