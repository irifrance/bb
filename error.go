// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package bb

import "errors"

var (
	OutOfBounds    = errors.New("Out of Bounds")
	BufferTooSmall = errors.New("Buffer To Small")
	NoReaderError  = errors.New("No Reader for Buffer")
	NoWriterError  = errors.New("No Writer for Buffer")
)
