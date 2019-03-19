package main

import "io"

// operation
type operation interface {

	// 自分自身の命令サイズを返す
	Size() int

	// ラベルの解決
	//
	// @param table --- ラベルテーブル
	//
	// @return エラー ラベルが見つからない場合
	Relocate(table map[string]int) error

	// オペレーションをバイナリとして出力
	Write(w io.Writer) (int, error)
}

// DB命令
type opDB struct {
	b []byte
}

func (o *opDB) Size() int {
	return len(o.b)
}

func (o *opDB) Relocate(_ map[string]int) error {
	return nil
}

func (o *opDB) Write(w io.Writer) (int, error) {
	return w.Write(o.b)
}

// RESB命令
type opRESB struct {
	size int
}

func (o *opRESB) Size() int {
	return o.size
}

func (o *opRESB) Relocate(_ map[string]int) error {
	return nil
}

func (o *opRESB) Write(w io.Writer) (int, error) {

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
