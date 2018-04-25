# bit buffers

Simple package for reading/writing bits from a memory backed []byte buffer

provides no wrappers around io.Reader/Writer and assumes the caller knows
that the backing buffer has sufficient size.
