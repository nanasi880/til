package main // import "github.com/nanasi880/til/os/day1"

import (
	"io/ioutil"
	"log"
)

// イメージのヘッダ
var imageHeader = []byte{
	0xEB, 0x4E, 0x90, 0x48, 0x45, 0x4C, 0x4C, 0x4F, 0x49, 0x50, 0x4C, 0x00, 0x02, 0x01, 0x01, 0x00,
	0x02, 0xE0, 0x00, 0x40, 0x0B, 0xF0, 0x09, 0x00, 0x12, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x40, 0x0B, 0x00, 0x00, 0x00, 0x00, 0x29, 0xFF, 0xFF, 0xFF, 0xFF, 0x48, 0x45, 0x4C, 0x4C, 0x4F,
	0x2D, 0x4F, 0x53, 0x20, 0x20, 0x20, 0x46, 0x41, 0x54, 0x31, 0x32, 0x20, 0x20, 0x20, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xB8, 0x00, 0x00, 0x8E, 0xD0, 0xBC, 0x00, 0x7C, 0x8E, 0xD8, 0x8E, 0xC0, 0xBE, 0x74, 0x7C, 0x8A,
	0x04, 0x83, 0xC6, 0x01, 0x3C, 0x00, 0x74, 0x09, 0xB4, 0x0E, 0xBB, 0x0F, 0x00, 0xCD, 0x10, 0xEB,
	0xEE, 0xF4, 0xEB, 0xFD, 0x0A, 0x0A, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0x77, 0x6F, 0x72,
	0x6C, 0x64, 0x0A,
}

func main() {

	bin := make([]byte, 0x168000)

	// ヘッダ付与
	copy(bin, imageHeader)

	bin[0x1FE] = 0x55
	bin[0x1FF] = 0xAA
	bin[0x200] = 0xF0
	bin[0x201] = 0xFF
	bin[0x202] = 0xFF

	bin[0x1400] = 0xF0
	bin[0x1401] = 0xFF
	bin[0x1402] = 0xFF

	if err := ioutil.WriteFile("hellos.img", bin, 0666); err != nil {
		log.Fatal(err)
	}
}
