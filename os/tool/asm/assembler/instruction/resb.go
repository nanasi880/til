package instruction

import "io"

// RESB命令
type RESB struct {
	size int
}

func NewRESB(size int) *RESB {
	return &RESB{
		size: size,
	}
}

func (o *RESB) Size() int {
	return o.size
}

func (o *RESB) Relocate(_ map[string]int) error {
	return nil
}

func (o *RESB) Write(w io.Writer) (int, error) {

	zero := make([]byte, 4096)
	size := o.size

	for size > 0 {

		if size < len(zero) {
			zero = zero[:size]
		}

		n, err := w.Write(zero)
		if err != nil {
			writtenSize := o.size - size + n
			return writtenSize, err
		}

		size -= n
	}

	return o.size, nil
}
