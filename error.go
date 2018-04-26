package bb

import "errors"

var (
	OutOfBounds    = errors.New("Out of Bounds")
	BufferTooSmall = errors.New("Buffer To Small")
)
