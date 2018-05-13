package bb

import "errors"

var (
	OutOfBounds    = errors.New("Out of Bounds")
	BufferTooSmall = errors.New("Buffer To Small")
	NoReaderError  = errors.New("No Reader for Buffer")
	NoWriterError  = errors.New("No Writer for Buffer")
)
