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
func (a *assembler) asm(sourceFile []byte) ([]byte, error) {

	// init
	reader := bufio.NewReader(bytes.NewReader(sourceFile))

	// 適当なサイズで1行分のデータを確保するためのバッファを作成
	line := make([]byte, 0, 1024)
	for {
		// データリセット
		line = line[:]

		a.sourceLineNumber += 1

	again:
		// 読めるところまで読む
		l, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			return a.relocate()
		}
		if err != nil {
			return nil, err
		}

		// 今回の読み込みで取得出来たデータは次回のReadLine()呼び出しの時点でスライスが書き換わるのでディープコーピーした上で
		// 後続データがあるなら引き続き読み込み
		line = append(line, l...)
		if isPrefix {
			goto again
		}

		// 1行分のデータを取得出来たのでそれを処理
		if err := a.line(line); err != nil {
			return nil, err
		}
	}
}

// アセンブリファイル1行分のデータの処理を開始
//
// @param line --- 1行分のデータ
//
// @return エラー
func (a *assembler) line(line []byte) error {

	var (
		tab        = []byte("\t")
		whiteSpace = []byte("    ")
	)

	// TAB文字は面倒なので空白に置換する
	line = bytes.ReplaceAll(line, tab, whiteSpace)

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

// オペレーションコード業をパースする
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

	// ニーモックに応じて処理を実施
	var (
		operations []operation
		err        error
	)
	switch mnemonic {
	case "DB":
		operations, err = a.mnemonicDB(parameter)

	default:
		return fmt.Errorf("error:%d unknown mnemonic `%s`", a.sourceLineNumber, mnemonic)
	}

	if err != nil {
		return fmt.Errorf("error:%d %s", a.sourceLineNumber, err.Error())
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

// DB命令
//
// @param parameter --- パラメーター
//
// @return オペレーション一覧、エラー
func (a *assembler) mnemonicDB(parameter string) ([]operation, error) {

	tokens, err := a.splitToken(parameter)
	if err != nil {
		return nil, err
	}

	var (
		result []operation
	)
	for _, tok := range tokens {

		switch tok := tok.(type) {

		case string:
			result = append(result, &opDB{
				b: []byte(tok),
			})

		case int:
			if tok < 0 || tok > 0xFF {
				return nil, fmt.Errorf("DB命令の即値は0x00 ~ 0xFFの範囲である必要がある")
			}
			result = append(result, &opDB{
				b: []byte{byte(tok)},
			})

		default:
			return nil, fmt.Errorf("internal %+v", tok)
		}
	}

	return result, nil
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

			if isString {
				result = append(result, string(token))
			} else {
				v, err := strconv.ParseInt(string(token), 0, 32)
				if err != nil {
					return nil, err
				}
				result = append(result, int(v))
			}
			isQuarto = false
			isString = false
			token = token[:]

		case ' ':
			if isQuarto {
				token = append(token, c)
				continue
			}

		default:
			token = append(token, c)
		}
	}

	return result, nil
}

func (a *assembler) relocate() ([]byte, error) {

	var (
		bin []byte
	)
	for _, o := range a.operations {

		if err := o.Relocate(a.labels); err != nil {
			return nil, err
		}

		bin = append(bin, o.Write()...)
	}

	return bin, nil
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
