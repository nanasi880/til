// Package lexer : 字句解析処理
package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/nanasi880/til/os/tool/asm/internal/runes"
)

type (
	Token string
	Line  []Token
	File  []Line
)

func (t Token) Last() uint8 {

	if len(t) == 0 {
		return 0
	}

	return t[len(t)-1]
}

// 字句解析実行
//
// @param src --- ソースコード
//
// @return 字句解析後のソースコード、エラー
func Analyze(src io.Reader) (File, error) {

	reader := bufio.NewReader(src)

	result := make(File, 0)

	// 適当なサイズで1行分のデータを確保するためのバッファを作成
	line := make([]byte, 0, 1024)
	for {
		// データリセット
		line = line[:0]

	again:
		// 読めるところまで読む
		l, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			return result, nil
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

		// 1行分のデータを解析
		if line, err := analyzeLine(bytesToRunes(line)); err == nil {
			if len(line) > 0 {
				result = append(result, line)
			}
		} else {
			return nil, err
		}
	}
}

// 行単位での字句解析実行
func analyzeLine(line []rune) (Line, error) {

	// TAB文字は面倒なので空白に置換する
	line = Clean(line)

	// 空行は無視
	if IsEmptyLine(line) {
		return nil, nil
	}

	tokens, err := SplitToken(line)
	if err != nil {
		return nil, err
	}

	return Line(tokens), nil
}

// Lexerが取り扱えるように行データをクリーンにする
// 引数で渡したlineには破壊的変更が行わなれる
//
// @param line --- １行分のテキストデータ
//
// @return クリーン後のデータ
func Clean(line []rune) []rune {

	// TAB文字は面倒なので空白に置換する
	line = ReplaceTab(line)

	// コメントを削除
	return TrimComment(line)
}

// 文字列をカンマ区切りのトークン列だと仮定して分割する
// ただし、最初のトークンは空白文字で区切られていると仮定される
// カンマから次のトークンまでの余分な空白は無視される
//
// この関数に渡す文字列はClean()でクリーニング済みである必要がある
//
// @param s --- 分割対象文字列
//
// @return トークンの一覧
func SplitToken(s []rune) ([]Token, error) {

	var result []Token

	// 余分な空白は捨てる
	s = runes.TrimSpace(s)

	// １つ目のトークンは空白で区切られているはず
	index := runes.Index(s, ' ')
	if index < 0 {
		// ラベルのようにそれ自体が１つのトークンとして完結している場合
		result = make([]Token, 0, 1)
		result = append(result, Token(s))
		return result, nil
	}

	result = append(result, Token(s[:index]))
	s = s[index:]

	// ２つ目以降のトークンはカンマで区切られているはず
	var (
		quotation bool
		token     = make([]rune, 0)
		prev      rune
	)
	escape := func() bool {
		return prev == '\\'
	}
	for i, c := range s {

		switch c {

		case '\\':
			if !quotation {
				token = append(token, c)
				goto next
			}

			if escape() {
				token = append(token, c)
				c = 0
			}

		case '"':
			if !escape() {
				quotation = !quotation
			}

			token = append(token, c)

		case ',':
			if escape() {
				return nil, fmt.Errorf("invalid escape: %d", i)
			}
			if quotation {
				token = append(token, c)
				goto next
			} else {
				if len(token) == 0 {
					return nil, fmt.Errorf("empty token: %d", i)
				}
				result = append(result, Token(token))
				quotation = false
				token = token[:0]
			}

		case ' ':
			if escape() {
				return nil, fmt.Errorf("invalid escape: %d", i)
			}
			if quotation {
				token = append(token, c)
			}

		default:
			if escape() {
				return nil, fmt.Errorf("invalid escape: %d", i)
			}
			token = append(token, c)
		}

	next:
		prev = c
	}

	// 最後まで読み切ったデータがあるならトークンとして処理する
	if len(token) > 0 {
		if quotation {
			// クォートが閉じられていない
			return nil, errors.New("quotation isn't closed")
		}
		result = append(result, Token(token))
	}

	return result, nil
}

// タブ文字を空白に置換する
// ただし、クォートされている部分はスキップする
//
// @param line --- 1行分のデータ
//
// @return タブを空白に置換した結果のデータ
func ReplaceTab(line []rune) []rune {

	var (
		quotation bool
		result    = make([]rune, 0, len(line))
		prev      rune
	)
	escape := func() bool {
		return prev == '\\'
	}
	for _, c := range line {

		switch c {

		case '\\':
			if !quotation {
				goto next
			}

			if escape() {
				result = append(result, c)
				c = 0
				goto next
			}

		case '"':
			if !escape() {
				quotation = !quotation
			}

		case '\t':
			if !quotation {
				c = ' '
			}
		}

		result = append(result, c)

	next:
		prev = c
	}

	return result
}

// コメントを除去する
//
// @param line --- 1行分のデータ
//
// @return コメントを除去した結果のデータ
func TrimComment(line []rune) []rune {

	var (
		quotation bool
		index     = -1
		prev      rune
	)
	escape := func() bool {
		return prev == '\\'
	}
	for i, c := range line {

		switch c {

		case '\\':
			if !quotation {
				goto next
			}

			if escape() {
				c = 0
				goto next
			}

		case '"':
			if !escape() {
				quotation = !quotation
			}

		case ';':
			if !quotation {
				index = i
				break
			}
		}

	next:
		prev = c
	}

	if index < 0 {
		return line
	}
	return line[:index]
}

// 事実上空行とみなせるかどうかを調べる
//
// この関数に渡す文字列はClean()でクリーニング済みである必要がある
//
// @param line --- 1行分のデータ
//
// @return 空行とみなせるかどうか 空白文字だけが存在するようなケースもtrueとみなす
func IsEmptyLine(line []rune) bool {

	for _, v := range line {
		if v != ' ' {
			return false
		}
	}
	return true
}

// バイトスライスをstringとして読み替える
// この関数で返却された文字列はunsafe経由でキャストされているため、元のバイトスライスと領域を共有していることに注意
//
// @param b --- バイト列
//
// @return 文字列
func bytesAsString(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// バイトスライスを一旦文字列として読み替えた後
// ルーンのスライスとしてallocし直して返却する
func bytesToRunes(b []byte) []rune {
	return []rune(bytesAsString(b))
}
