package assembler

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"go.nanasi880.dev/rpn"

	"github.com/nanasi880/til/os/tool/asm/assembler/instruction"
	"github.com/nanasi880/til/os/tool/asm/assembler/lexer"
	"github.com/nanasi880/til/os/tool/asm/internal"
)

func (a *Assembler) parseMnemonic(mnemonic lexer.Token, parameters []lexer.Token) error {

	var (
		err error
	)
	switch mnemonic {

	// data byte
	case "DB":
		err = a.mnemonicDB(parameters)

	// data word
	case "DW":
		err = a.mnemonicMultiWord(parameters, 2)

	// data double word
	case "DD":
		err = a.mnemonicMultiWord(parameters, 4)

	// reserve byte
	case "RESB":
		err = a.mnemonicRESB(parameters)

	default:
		return fmt.Errorf("error:%d unknown mnemonic `%s`", a.sourceLineNumber, mnemonic)
	}

	if err != nil {
		return fmt.Errorf("error:%d %s", a.sourceLineNumber, err.Error())
	}

	return nil
}

// 変数解決のリゾルバを取得する
//
// @return リゾルバ
func (a *Assembler) Resolver() rpn.Resolver {

	return func(name string) (decimal.Decimal, error) {

		if name == "$" {
			return decimal.New(a.address, 0), nil
		}
		return decimal.Zero, fmt.Errorf("undeclared variable: %s", name)
	}
}

// トークン列をパラメーターだと仮定してデコードする
// トークンがクォートされていない場合、それを式と解釈する
// それ以外はレジスタ名も含めてstringとして取り扱う
//
// @param parameters --- 分割対象文字列
//
// @return *rpn.RPN or stringの混合スライス、エラー
func (a *Assembler) decodeParameters(parameters []lexer.Token) ([]interface{}, error) {

	// トークンがエスケープされているかどうかを返す
	quoted := func(tok lexer.Token) bool {
		if tok[0] == '"' && tok.Last() == '"' && len(tok) >= 2 {
			return true
		}
		return false
	}

	var result []interface{}
	for _, p := range parameters {

		if quoted(p) {
			result = append(result, string(p))
		} else {
			rpnObject, err := rpn.Parse(string(p))
			if err != nil {
				return nil, err
			}
			result = append(result, rpnObject)
		}
	}

	return result, nil
}

// DB命令
//
// @param parameters --- パラメーター
//
// @return オペレーション一覧、エラー
func (a *Assembler) mnemonicDB(parameters []lexer.Token) error {

	return a.mnemonicMultiWordWithConverter(parameters, func(v interface{}) ([]byte, error) {

		switch v := v.(type) {

		case string:
			v = strings.TrimFunc(v, func(r rune) bool {
				return r == '"'
			})
			return []byte(v), nil

		case int64:
			if v > 0xFF || v < 0 {
				return nil, fmt.Errorf("DB命令の即値は0x00 ~ 0xFFの範囲である必要がある")
			}
			return []byte{byte(v)}, nil

		default:
			return nil, nil
		}
	})
}

// DW命令 / DD命令
//
// @param parameters --- パラメーター
// @param size       --- 命令サイズ
//                       2ならDW、4ならDDと解釈される
//
// @return オペレーション一覧、エラー
func (a *Assembler) mnemonicMultiWord(parameters []lexer.Token, size int) error {

	var (
		mnemonic       = "DW"
		max      int64 = 0xFFFF
	)
	if size == 4 {
		mnemonic = "DD"
		max = 0xFFFFFFFF
	}
	return a.mnemonicMultiWordWithConverter(parameters, func(v interface{}) ([]byte, error) {

		i, ok := v.(int64)
		if !ok {
			return nil, fmt.Errorf("%s命令は文字列は使用できない", mnemonic)
		}

		if i > max || i < 0 {
			return nil, fmt.Errorf("%s命令の即値は0x00 ~ 0x%Xの範囲である必要がある", mnemonic, max)
		}

		if size == 4 {
			bytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bytes, uint32(i))
			return bytes, nil
		} else {
			bytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(bytes, uint16(i))
			return bytes, nil
		}
	})
}

func (a *Assembler) mnemonicMultiWordWithConverter(parameters []lexer.Token, c func(v interface{}) ([]byte, error)) error {

	if len(parameters) == 0 {
		return fmt.Errorf("最低1つのパラメーターが必要")
	}

	decodedParameters, err := a.decodeParameters(parameters)
	if err != nil {
		return err
	}
	for _, p := range decodedParameters {

		switch p := p.(type) {

		case string:
			b, err := c(p)
			if err != nil {
				return err
			}
			a.mnemonics = append(a.mnemonics, instruction.NewDB(b))
			a.address += int64(len(b))

		case *rpn.RPN:
			d, err := p.Eval(a.Resolver())
			if err != nil {
				return err
			}

			b, err := c(d.IntPart())
			if err != nil {
				return err
			}
			a.mnemonics = append(a.mnemonics, instruction.NewDB(b))
			a.address += int64(len(b))

		default:
			return fmt.Errorf("internal: %#v", p)
		}
	}

	return nil
}

// RESB命令
//
// @param parameter --- パラメーター
//
// @return エラー
func (a *Assembler) mnemonicRESB(parameters []lexer.Token) error {

	if len(parameters) != 1 {
		return fmt.Errorf("RESB命令は1つのパラメーターが必要")
	}

	rpnObject, err := rpn.Parse(string(parameters[0]))
	if err != nil {
		return err
	}
	d, err := rpnObject.Eval(a.Resolver())
	if err != nil {
		return err
	}
	v := d.IntPart()

	if v < 0 {
		return fmt.Errorf("RESB underflow: %d", v)
	}
	if v > internal.MaxInt {
		return fmt.Errorf("RESB overflow: %d", v)
	}

	resb := instruction.NewRESB(v)
	a.mnemonics = append(a.mnemonics, resb)
	a.address += resb.Size()

	return nil
}
