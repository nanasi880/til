package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type assembler struct {
	origin           int            // 命令配置基準位置 ORG命令でセットされる
	address          int            // originから現在の命令位置のオフセット
	sourceLineNumber int            // 現在解析しているソースコードの行番号
	labels           map[string]int // ラベルの名前:addressの対応表
	operations       []operation    // バイナリ先頭からのオペコード一覧
}

// 指定したファイルのアセンブルを開始
func (a *assembler) asm(sourceFile io.Reader, out io.Writer) error {

	// init
	reader := bufio.NewReader(sourceFile)

	// 適当なサイズで1行分のデータを確保するためのバッファを作成
	line := make([]byte, 0, 1024)
	for {
		// データリセット
		line = line[:0]

		a.sourceLineNumber += 1

	again:
		// 読めるところまで読む
		l, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			return a.relocate(out)
		}
		if err != nil {
			return err
		}

		// 今回の読み込みで取得出来たデータは次回のReadLine()呼び出しの時点でスライスが書き換わるのでディープコーピーした上で
		// 後続データがあるなら引き続き読み込み
		line = append(line, l...)
		if isPrefix {
			goto again
		}

		// 1行分のデータを取得出来たのでそれを処理
		if err := a.line(line); err != nil {
			return err
		}
	}
}

// アセンブリファイル1行分のデータの処理を開始
//
// @param line --- 1行分のデータ
//
// @return エラー
func (a *assembler) line(line []byte) error {

	// TAB文字は面倒なので空白に置換する
	line = a.replaceTab(line)

	// コメントより後ろは削除
	line = a.trimComment(line)

	// 空行は無視
	if a.isEmpty(line) {
		return nil
	}

	if line[0] != ' ' {

		// 先頭に余白が無い場合、それはコメント行 or ラベル
		if err := a.parseLabel(line); err != nil {
			return err
		}

	} else {

		// それ以外は命令行もしくは空行
		if err := a.parseOpCode(line); err != nil {
			return err
		}
	}

	return nil
}

// ラベル行をパースする
func (a *assembler) parseLabel(line []byte) error {

	// ラベル行は必ずコロンで終端しているはず
	index := bytes.IndexByte(line, ':')
	if index < 0 {
		return fmt.Errorf("error:%d ラベル名がコロンで終端していない", a.sourceLineNumber)
	}

	// すでにラベル名が存在しているのはコンパイルエラー
	label := string(line[:index])
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
func (a *assembler) parseOpCode(line []byte) error {

	// オペコード解析に空白は邪魔なので削除してしまう
	line = bytes.TrimSpace(line)

	// 最初の空白までを取り出し、その結果がニーモックのはず
	var (
		mnemonic  string
		parameter string
	)
	if index := bytes.IndexByte(line, ' '); index > 0 {
		mnemonic = string(line[:index])
		parameter = string(bytes.TrimSpace(line[index:]))
	} else {
		mnemonic = string(line)
	}
	mnemonic = strings.ToUpper(mnemonic)

	operations, err := a.parseMnemonic(mnemonic, parameter)
	if err != nil {
		return err
	}

	// 命令サイズ分だけアドレスをオフセットし、命令一覧を結合
	var operationSize int
	for _, o := range operations {
		operationSize += o.Size()
	}
	a.address += operationSize
	a.operations = append(a.operations, operations...)

	return nil
}

// 文字列をカンマ区切りのトークン列だと仮定して分割する
// 0xで始まるテキストでかつそれがクオートされていない場合、中身をintと仮定しパースする
// それ以外はレジスタ名も含めてstringとして取り扱う
// カンマから次のトークンまでの余分な空白は無視される
//
// @param s --- 分割対象文字列
//
// @return int or stringの混合スライス、エラー
func (a *assembler) splitToken(s string) ([]interface{}, error) {

	var (
		result   []interface{}
		isQuarto bool
		token    = make([]byte, 0)
		isString bool
	)

	doParse := func() error {

		if isString {
			result = append(result, string(token))
		} else {
			v, err := strconv.ParseUint(string(token), 0, 32)
			if err != nil {
				return err
			}
			result = append(result, v)
		}
		return nil
	}

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch c {

		case '"':
			isQuarto = !isQuarto
			isString = true

		case ',':
			if isQuarto {
				token = append(token, c)
				continue
			}

			if err := doParse(); err != nil {
				return nil, err
			}

			isQuarto = false
			isString = false
			token = token[:0]

		case ' ':
			if isQuarto {
				token = append(token, c)
				continue
			}

		default:
			token = append(token, c)
		}
	}

	// 最後まで読み切ったデータも解析する
	if len(token) > 0 {
		if err := doParse(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (a *assembler) relocate(out io.Writer) error {

	bo := bufio.NewWriter(out)

	for _, o := range a.operations {

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

// 事実上空行とみなせるかどうかを調べる
//
// @param line --- 1行分のデータ
//
// @return 空行とみなせるかどうか 空白文字だけが存在するようなケースもtrueとみなす
func (a *assembler) isEmpty(line []byte) bool {
	for _, v := range line {
		if v != ' ' {
			return false
		}
	}
	return true
}

// タブ文字を空白に置換する
// ただし、クォートされている部分はスキップする
//
// @param line --- 1行分のデータ
//
// @return タブを空白に置換した結果のデータ
func (a *assembler) replaceTab(line []byte) []byte {

	var (
		isQuote bool
		result = make([]byte, 0, len(line))
	)
	for _, c := range line {

		switch c {
		case '"':
			isQuote = !isQuote

		case '\t':
			if !isQuote {
				c = ' '
			}
		}

		result = append(result, c)
	}

	return result
}

// コメントを除去する
//
// @param line --- 1行分のデータ
//
// @return コメントを除去した結果のデータ
func (a *assembler) trimComment(line []byte) []byte {

	var (
		isQuote bool
		index   = -1
	)
	for i, c := range line {

		switch c {
		case '"':
			isQuote = !isQuote

		case ';':
			if !isQuote {
				index = i
				break
			}
		}
	}

	if index < 0 {
		return line
	}
	return line[:index]
}
