package instruction

import "io"

// DB命令
type DB struct {
	b []byte
}

func NewDB(b []byte) *DB {
	return &DB{
		b: b,
	}
}

func (o *DB) Size() int {
	return len(o.b)
}

func (o *DB) Relocate(_ map[string]int) error {
	return nil
}

func (o *DB) Write(w io.Writer) (int, error) {
	return w.Write(o.b)
}
