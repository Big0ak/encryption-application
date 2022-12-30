package main

import (
	// "bytes"
	// "encoding/binary"

	// "log"
	// "os"
	// "strconv"

	"github.com/Big0ak/AES/AES"
)

func main() {


	//var word_key string = "qwertyuiasdfghjk"
	key_byte := [16]byte{ 0x2b, 0x7e, 0x15, 0x16,
						0x28, 0xae, 0xd2, 0xa6,
						0xab, 0xf7, 0x15, 0x88,
						0x09, 0xcf, 0x4f, 0x3c}
	
	text_byte := [16]byte{ 0x32, 0x88, 0x31, 0xe0,
						0x43, 0x5a, 0x31, 0x37,
						0xf6, 0x30, 0x98, 0x07,
						0xa8, 0x8d, 0xa2, 0x34 };	

	AES.Encrypt(key_byte[:], text_byte[:])

	// file, err := os.Open("text.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// byte := make([]byte, 16)
	// kol, err := file.Read(byte)
	// fmt.Print(kol, "\n", byte)

	// buf := bytes.NewReader(byte)
	// var bin float64
	// err = binary.Read(buf, binary.BigEndian, &bin)
	// fmt.Print(bin, "\n")

	// fi, err := file.Stat()
	// fmt.Print(fi.Size())
	// file.Close()
	
	// /////////////////////////////////////////////
	// file, err = os.Create("binary.bin")
	// file.WriteString(strconv.FormatFloat(bin, 'f', 6, 64))
	// file.Close()
	// ///////////////////////////////////////

	// buf2 := new(bytes.Buffer)
	// err = binary.Write(buf2, binary.BigEndian, bin)
	// if err != nil {
	// 	log.Fatal("binary.Write failed:", err)
	// }
	// fmt.Print(buf2.Bytes())

	// file, err = os.Create("test.txt")
	// kol, err = file.Write(buf2.Bytes())
	// file.Close()
}
