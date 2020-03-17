package instruction

import (
	"io"

	"go.nanasi880.dev/rpn"
)

// DB命令
type DB struct {
	b   []byte
	rpn *rpn.RPN
}

func NewDB(b []byte) *DB {
	return &DB{
		b: b,
	}
}

func (o *DB) Size() int64 {
	return int64(len(o.b))
}

func (o *DB) Relocate(table map[string]int64) error {
	return nil
}

func (o *DB) Write(w io.Writer) (int64, error) {
	n, err := w.Write(o.b)
	return int64(n), err
}
