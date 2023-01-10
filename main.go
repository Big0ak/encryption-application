// Oficial document FIPS 197, Advanced Encryption Standard (AES)
// https://csrc.nist.gov/csrc/media/publications/fips/197/final/documents/fips-197.pdf

// На русском
// http://crypto.pp.ua/wp-content/uploads/2010/03/aes.pdf

// Пример реализации
// https://programmer.group/c-implementation-of-aes-encryption-algorithms.html

package main

import (
	"io"
	"log"
	"os"
	"time"

	//"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/Big0ak/AES/AES"
)

const (
	ExpCrypto   = ".crypto"
	CipherFile  = "encrypted" + ExpCrypto
	DectyptFile = "decrypt"
)

var sourceFileLoaded = false // сигнал о том, что файл загружен
var endWork = false          // сигнал о том, что шифрование/дешифрование прошло над файлом успешно => финальное окно

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

	var (
		key        = ""
		readerFile fyne.URIReadCloser
	)

	a := app.New()
	w := a.NewWindow(" ")
	//a.Settings().SetTheme(theme.DarkTheme())
	w.Resize(fyne.NewSize(300, 400))
	w.CenterOnScreen()   // окно по центру экрана
	w.SetFixedSize(true) // нельзя менять размер
	w.SetMaster()        // главное окно

	//-------------------------- ЭКРАН ЗАГРУЗКИ ФАЙЛА --------------------------

	img_fileUpload := canvas.NewImageFromFile("file_upload.png")
	img_fileUpload.Resize(fyne.NewSize(80, 80))
	img_fileUpload.Move(fyne.NewPos(100, 160))

	btn_openFile := widget.NewButton(" ", func() {
		w2 := a.NewWindow(" ")
		w2.Resize(fyne.NewSize(525, 370))
		w2.CenterOnScreen()
		w2.SetFixedSize(true)
		dialog.ShowFileOpen(
			func(uc fyne.URIReadCloser, err error) {
				readerFile = uc
				sourceFileLoaded = true
				w2.Close()
			},
			w2,
		)
		w2.Show()
	})
	btn_openFile.Resize(fyne.NewSize(290, 100))
	btn_openFile.Move(fyne.NewPos(0, 150))

	img_textFooter := canvas.NewImageFromFile("text_footer.png")
	img_textFooter.Resize(fyne.NewSize(300, 80))
	img_textFooter.Move(fyne.NewPos(0, 285))

	cont_UploadFile := container.NewWithoutLayout(
		btn_openFile,
		img_fileUpload,
		img_textFooter,
	)

	//-------------------------- ЭКРАН ВВОДА КЛЮЧА --------------------------

	img_fileIcon := canvas.NewImageFromFile("file_icon.png")
	box_fileIcon := container.NewHBox(
		layout.NewSpacer(),
		container.New(
			layout.NewGridWrapLayout(fyne.NewSize(96, 96)),
			img_fileIcon,
		),
		layout.NewSpacer(),
	)

	nameFile := binding.NewString()
	box_nameFile := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabelWithData(nameFile),
		layout.NewSpacer(),
	)

	wid_inputKey := widget.NewPasswordEntry()
	wid_inputKey.SetPlaceHolder("Введите ключ")

	btn_crypto := widget.NewButton("Crypto", func() {
		// TODO: валидация ключа
		key = wid_inputKey.Text
		if readerFile.URI().Extension() == ExpCrypto {
			// Расшифровка
			enc, err := io.ReadAll(readerFile)
			if err != nil {
				log.Fatal(err)
			}
			plain_byte, ext, err := AES.Decrypt([]byte(key), enc[:])
			if err == nil {
				file, _ := os.Create(DectyptFile + ext)
				file.Write(plain_byte)
				file.Close()
			} else {
				log.Fatal(err)
			}
		} else {
			// Шифрование
			plain, err := io.ReadAll(readerFile)
			if err != nil {
				log.Fatal(err)
			}
			enc_byte, err := AES.Encrypt([]byte(key), plain[:], []byte(readerFile.URI().Extension()))
			if err == nil {
				file, _ := os.Create(CipherFile)
				file.Write(enc_byte)
				file.Close()
			} else {
				log.Fatal(err)
			}
		}

		endWork = true
	})

	cont_enteringKey := container.NewVBox(
		box_fileIcon,
		box_nameFile,
		wid_inputKey,
		btn_crypto,
	)
	cont_enteringKey.Hide()

	//-------------------------- ЭКРАН ЗАВЕРШЕНИЯ РАБОТЫ --------------------------

	wid_endSuccess := widget.NewLabel("Работы на файлом завершена")

	cont_endScene := container.NewVBox(
		wid_endSuccess,
	)
	cont_endScene.Hide()

	//-----------------------------------------------------------------------------

	// Общий контейнер
	cont := container.NewVBox(
		cont_UploadFile,
		cont_enteringKey,
		cont_endScene,
	)

	w.SetContent(cont)
	w.Show()

	// Обновление экранов
	go func() {
		for range time.Tick(time.Second) {
			if sourceFileLoaded {
				cont_UploadFile.Hide()
				nameFile.Set(readerFile.URI().Name())
				cont_enteringKey.Show()
				sourceFileLoaded = false
			}

			if endWork {
				cont_enteringKey.Hide()
				cont_endScene.Show()
				endWork = false
			}
		}
	}()

	a.Run()
}
