package assembler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nanasi880/til/os/tool/asm/assembler/instruction"
	"github.com/nanasi880/til/os/tool/asm/assembler/lexer"
)

type Assembler struct {
	origin           int                    // 命令配置基準位置 ORG命令でセットされる
	address          int                    // originから現在の命令位置のオフセット
	sourceLineNumber int                    // 現在解析しているソースコードの行番号
	labels           map[string]int         // ラベルの名前:addressの対応表
	mnemonics        []instruction.Mnemonic // バイナリ先頭からのオペコード一覧
}

// 新しいアセンブラインスタンスを作成
func New() *Assembler {
	return new(Assembler)
}

// 指定したファイルのアセンブルを開始
func (a *Assembler) Exec(sourceFile io.Reader, out io.Writer) error {

	file, err := lexer.Analyze(sourceFile)
	if err != nil {
		return err
	}

	a.sourceLineNumber = 1
	for _, line := range file {
		if err := a.line(line); err != nil {
			return err
		}
		a.sourceLineNumber += 1
	}

	for _, m := range a.mnemonics {
		if err := m.Relocate(a.labels); err != nil {
			return err
		}
	}

	w := bufio.NewWriter(out)
	for _, m := range a.mnemonics {
		if _, err := m.Write(w); err != nil {
			return err
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

// アセンブリファイル1行分のデータの処理を開始
//
// @param line --- 1行分のデータ
//
// @return エラー
func (a *Assembler) line(line lexer.Line) error {

	if line[0].Last() == ':' {
		return a.parseLabel(line)
	} else {
		return a.parseOpCode(line)
	}
}

// ラベル行をパースする
func (a *Assembler) parseLabel(line lexer.Line) error {

	// 末尾のコロンを削除
	label := string(line[0])
	label = label[:len(label)-1]

	// 既にラベル名が存在しているのはコンパイルエラー
	if _, ok := a.labels[label]; ok {
		return fmt.Errorf("error:%d ラベル名 %s は既に使用されています", a.sourceLineNumber, label)
	}

	// ラベル名と現在のオフセットアドレスを記憶
	if a.labels == nil {
		a.labels = make(map[string]int)
	}
	a.labels[label] = a.address

	return nil
}

// オペレーションコード行をパースする
func (a *Assembler) parseOpCode(line lexer.Line) error {

	var (
		mnemonic   = line[0]
		parameters = line[1:]
	)
	err := a.parseMnemonic(mnemonic, parameters)
	if err != nil {
		return err
	}

	return nil
}

func (a *Assembler) relocate(out io.Writer) error {

	bo := bufio.NewWriter(out)

	for _, o := range a.mnemonics {

		if err := o.Relocate(a.labels); err != nil {
			return err
		}

		_, err := o.Write(bo)
		if err != nil {
			return err
		}
	}

	if err := bo.Flush(); err != nil {
		return err
	}

	return nil
}
