// Oficial document FIPS 197, Advanced Encryption Standard (AES)
// https://csrc.nist.gov/csrc/media/publications/fips/197/final/documents/fips-197.pdf

// На русском
// http://crypto.pp.ua/wp-content/uploads/2010/03/aes.pdf

// Пример реализации
// https://programmer.group/c-implementation-of-aes-encryption-algorithms.html

package main

import (
	"log"
	"os"
	"path/filepath"

	//"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Big0ak/AES/AES"
)

const ExpCrypto = ".crypto"

var SourceFileLoaded = false // сигнал о том, что фай загружен
var EndWork = false // сигнал о том, что шифрование/дешифрование прошло на файлом успешно => финальное окно


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

	// key := "4m5n7q8r9t2j3k4n"
	// nameFile := "test.jpg"
	// cipherFile := "chifer.aes"
	// dectyptFile := "decrypt"

	// plain, err := os.ReadFile(nameFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// enc_byte, err := AES.Encrypt([]byte(key), plain[:], []byte(filepath.Ext("/" + nameFile)))
	// if err == nil {
	// 	file, _ := os.Create(cipherFile)
	// 	file.Write(enc_byte)
	// 	file.Close()
	// } else {
	// 	log.Fatal(err)
	// }

	// // ----------------------------------------------------------------------------------

	// enc, err := os.ReadFile(cipherFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// plain_byte, ext, err := AES.Decrypt([]byte(key), enc[:])
	// if err == nil {
	// 	file, _ := os.Create(dectyptFile+ext)
	// 	file.Write(plain_byte)
	// 	file.Close()
	// } else {
	// 	log.Fatal(err)
	// }

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow(" ")
	w.Resize(fyne.NewSize(300, 400))
	w.CenterOnScreen() // окно по центру экрана
	w.SetFixedSize(true) // нельзя менять размер
	w.SetMaster() // главное окно

	var ( 
		sourceFile = ""
		cipherFile = "chifer" + ExpCrypto
		dectyptFile = "decrypt"
		key = ""
	)

	openFile := widget.NewButtonWithIcon("Open file", theme.FileIcon(), func() {
		w2 := a.NewWindow(" ")
		w2.Resize(fyne.NewSize(525,370))
		w2.CenterOnScreen()
		w2.SetFixedSize(true)
		dialog.ShowFileOpen(
			func(uc fyne.URIReadCloser, err error) {
				sourceFile = uc.URI().Name()
				SourceFileLoaded = true
				w2.Close()
			},
			w2,
		)
		w2.Show()
		
	})
	openFile.Resize(fyne.NewSize(300, 250))
	openFile.Move(fyne.NewPos(0,90))

	uploadFile := container.NewWithoutLayout(
		openFile,
	)
	
/////////////////////////////////////////////////////////////////////

	nameFile := binding.NewString()
	nameFileWid := widget.NewLabelWithData(nameFile)

	inputKey := widget.NewEntry()
	inputKey.SetPlaceHolder("Введите ключ")
	
	btn := widget.NewButton("Crypto", func() {
		key = inputKey.Text

		if (filepath.Ext("/" + sourceFile)) == ExpCrypto{
			enc, err := os.ReadFile(sourceFile)
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
		} else {
			plain, err := os.ReadFile(sourceFile)
			if err != nil {
				log.Fatal(err)
			}
			enc_byte, err := AES.Encrypt([]byte(key), plain[:], []byte(filepath.Ext("/" + sourceFile)))
			if err == nil {
				file, _ := os.Create(cipherFile)
				file.Write(enc_byte)
				file.Close()
			} else {
				log.Fatal(err)
			}
		}

		EndWork = true
	})

	enteringKey := container.NewVBox(
		nameFileWid,
		inputKey,
		btn,
	)
	enteringKey.Hide()
//////////////////////////////////////////////////////////////////

	endSuccess := widget.NewLabel("Работы на файлом завершена")

	endScene := container.NewVBox(
		endSuccess,
	)

	endScene.Hide()
/////////////////////////////////////////////////////////////////

	// Общий контейнер
	cont := container.NewVBox(
		uploadFile,
		enteringKey,
		endScene,
	)

	w.SetContent(cont)
	w.Show()

	go func() {
		for ; ; {
			if SourceFileLoaded {
				uploadFile.Hide()
				nameFile.Set(sourceFile)
				enteringKey.Show()
				SourceFileLoaded = false
			}

			if EndWork {
				enteringKey.Hide()
				endScene.Show()
				EndWork = false
			}
		}
	} ()

	a.Run()
}

// ic, _ := fyne.LoadResourceFromPath("123.jpg")