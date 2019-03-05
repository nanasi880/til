package main

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
	Write() []byte
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

func (o *opDB) Write() []byte {
	return o.b
}
