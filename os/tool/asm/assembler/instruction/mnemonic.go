package instruction

import "io"

// Mnemonic
type Mnemonic interface {

	// 自分自身の命令サイズを返す
	Size() int64

	// ラベルの解決
	//
	// @param table --- ラベルテーブル
	//
	// @return エラー ラベルが見つからない場合
	Relocate(table map[string]int64) error

	// オペレーションをバイナリとして出力
	Write(w io.Writer) (int64, error)
}
