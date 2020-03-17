package instruction

import "io"

// RESB命令
type RESB struct {
	size int64
}

func NewRESB(size int64) *RESB {
	return &RESB{
		size: size,
	}
}

func (o *RESB) Size() int64 {
	return o.size
}

func (o *RESB) Relocate(table map[string]int64) error {
	return nil
}

func (o *RESB) Write(w io.Writer) (int64, error) {

	zero := make([]byte, 4096)
	size := o.size

	for size > 0 {

		if size < int64(len(zero)) {
			zero = zero[:size]
		}

		n, err := w.Write(zero)
		if err != nil {
			writtenSize := o.size - size + int64(n)
			return writtenSize, err
		}

		size -= int64(n)
	}

	return o.size, nil
}
