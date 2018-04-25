# bit buffers

Simple package for reading/writing bits from a memory backed []byte buffer
bb:

- provides no wrappers around io.Reader/Writer.
- has no errors in return values
- grows the backing []byte as needed and capable
- panics on Read* beyond last byte in the backing buffer
