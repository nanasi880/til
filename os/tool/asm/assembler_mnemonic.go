package main

import (
	"fmt"
	"strconv"

	"github.com/nanasi880/til/os/tool/asm/internal"
)

func (a *assembler) parseMnemonic(mnemonic string, parameter string) ([]operation, error) {

	var (
		operations []operation
		err        error
	)
	switch mnemonic {

	// data byte
	case "DB":
		operations, err = a.mnemonicDB(parameter)

	// reserve byte
	case "RESB":
		operations, err = a.mnemonicRESB(parameter)

	default:
		return nil, fmt.Errorf("error:%d unknown mnemonic `%s`", a.sourceLineNumber, mnemonic)
	}

	if err != nil {
		return nil, fmt.Errorf("error:%d %s", a.sourceLineNumber, err.Error())
	}

	return operations, nil
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

// RESB命令
//
// @param parameter --- パラメーター
//
// @return オペレーション一覧、エラー
func (a *assembler) mnemonicRESB(parameter string) ([]operation, error) {

	size, err := strconv.ParseInt(parameter, 0, internal.IntSize)
	if err != nil {
		return nil, err
	}

	result := make([]operation, 0, 1)
	result = append(result, &opRESB{
		size: int(size),
	})

	return result, nil
}
