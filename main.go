package main

import (
	"log"
	"os"
	"path/filepath"
	//"fmt"

	"github.com/Big0ak/AES/AES"
)

func main() {
	// key_byte := [16]byte{ 0x2b, 0x7e, 0x15, 0x16,
	// 					0x28, 0xae, 0xd2, 0xa6,
	// 					0xab, 0xf7, 0x15, 0x88,
	// 					0x09, 0xcf, 0x4f, 0x3c}
	// text_byte := [16]byte{ 0x32, 0x88, 0x31, 0xe0,
	// 					0x43, 0x5a, 0x31, 0x37,
	// 					0xf6, 0x30, 0x98, 0x07,
	// 					0xa8, 0x8d, 0xa2, 0x34 };	
	// enc_byte := AES.Encrypt(key_byte[:], text_byte[:])
	// plain_byte := AES.Decrypt(key_byte[:], enc_byte[:])
	// fmt.Print(plain_byte)

	key := "4m5n7q8r9t2j3k4n"
	nameFile := "test.jpg"
	cipherFile := "chifer.aes"
	dectyptFile := "decrypt"

	plain, err := os.ReadFile(nameFile)
	if err != nil {
		log.Fatal(err)
	}
	enc_byte, err := AES.Encrypt([]byte(key), plain[:], []byte(filepath.Ext("/" + nameFile)))
	if err == nil {
		file, _ := os.Create(cipherFile)
		file.Write(enc_byte)
		file.Close()
	} else {
		log.Fatal(err)
	}

	// ----------------------------------------------------------------------------------

	enc, err := os.ReadFile(cipherFile)
	if err != nil {
		log.Fatal(err)
	}
	plain_byte, ext, err := AES.Decrypt([]byte(key), enc[:])
	if err == nil {
		file, _ := os.Create(dectyptFile+ext)
		file.Write(plain_byte)
		file.Close()
	} else {
		log.Fatal(err)
	}
}