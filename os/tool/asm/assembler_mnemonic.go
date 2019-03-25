package main

import (
	"encoding/binary"
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

	// data word
	case "DW":
		operations, err = a.mnemonicMultiWord(parameter, 2)

	// data double word
	case "DD":
		operations, err = a.mnemonicMultiWord(parameter, 4)

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

		case uint64:
			if tok > 0xFF {
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

// DW命令 / DD命令
//
// @param parameter --- パラメーター
// @param size      --- 命令サイズ
//                      2ならDW、4ならDDと解釈される
//
// @return オペレーション一覧、エラー
func (a *assembler) mnemonicMultiWord(parameter string, size int) ([]operation, error) {

	if size == 2 {
		return a.mnemonicMultiWordWithConverter(parameter, func(u uint64) ([]byte, error) {

			if u > 0xFFFF {
				return nil, fmt.Errorf("DW命令の即値は0x0000 ~ 0xFFFFの範囲である必要がある")
			}

			bytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(bytes, uint16(u))

			return bytes, nil
		})
	}

	return a.mnemonicMultiWordWithConverter(parameter, func(u uint64) ([]byte, error) {

		if u > 0xFFFFFFFF {
			return nil, fmt.Errorf("DD命令の即値は0x00000000 ~ 0xFFFFFFFFの範囲である必要がある")
		}

		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, uint32(u))

		return bytes, nil
	})
}

func (a *assembler) mnemonicMultiWordWithConverter(parameter string, c func(uint64) ([]byte, error)) ([]operation, error) {

	tokens, err := a.splitToken(parameter)
	if err != nil {
		return nil, err
	}

	var (
		result []operation
	)

	for _, tok := range tokens {

		switch tok := tok.(type) {

		case uint64:

			bytes, err := c(tok)
			if err != nil {
				return nil, err
			}

			result = append(result, &opDB{
				b: bytes,
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
