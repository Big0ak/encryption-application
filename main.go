// Oficial document FIPS 197, Advanced Encryption Standard (AES)
// https://csrc.nist.gov/csrc/media/publications/fips/197/final/documents/fips-197.pdf

// На русском
// http://crypto.pp.ua/wp-content/uploads/2010/03/aes.pdf

// Пример реализации
// https://programmer.group/c-implementation-of-aes-encryption-algorithms.html

package main

import (
	"errors"
	"image/color"
	"io"
	"log"
	"os"
	"time"

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
	var (
		key        = ""
		readerFile fyne.URIReadCloser
	)

	a := app.New()
	w := a.NewWindow(" ")
	w.Resize(fyne.NewSize(300, 400))
	w.CenterOnScreen()   // окно по центру экрана
	w.SetFixedSize(true) // нельзя менять размер
	w.SetMaster()        // главное окно
	icon, _ := fyne.LoadResourceFromPath("icon.png")
	w.SetIcon(icon)

	//-------------------------- ЭКРАН ЗАГРУЗКИ ФАЙЛА --------------------------

	text_title := canvas.NewText("Шифрование файлов AES", color.Black)
	text_title.Alignment = fyne.TextAlignCenter
	text_title.TextSize = 16
	text_title.TextStyle.Bold = true
	text_title.Resize(fyne.NewSize(280,80))
	text_title.Move(fyne.NewPos(0,20))

	img_fileUpload := canvas.NewImageFromFile("file_upload.png")
	img_fileUpload.Resize(fyne.NewSize(80, 80))
	img_fileUpload.Move(fyne.NewPos(100, 160))

	btn_openFile := container.New(
		layout.NewMaxLayout(),

		widget.NewButton(" ", func() {
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
		}),
		canvas.NewHorizontalGradient(color.RGBA{12, 0, 255, 255}, color.Transparent),
		
	)
	btn_openFile.Resize(fyne.NewSize(290, 100))
	btn_openFile.Move(fyne.NewPos(0, 150))

	img_textFooter := canvas.NewImageFromFile("text_footer.png")
	img_textFooter.Resize(fyne.NewSize(300, 80))
	img_textFooter.Move(fyne.NewPos(0, 285))

	cont_UploadFile := container.NewWithoutLayout(
		text_title,
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
	wid_inputKey.Validator = func (input string) error {
		if len(input) != 16 {
			return errors.New("")
		}
		return nil
	}

	checkKey := binding.NewString()
	box_checkKey := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabelWithData(checkKey),
		layout.NewSpacer(),
	)

	btn_crypto := widget.NewButton("Crypto", func() {
		if err := wid_inputKey.Validate(); err == nil{

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
		} else {
			checkKey.Set("ключ минимум 16 символов")
		}
	})

	cont_enteringKey := container.NewVBox(
		box_fileIcon,
		box_nameFile,
		box_checkKey,
		wid_inputKey,
		btn_crypto,
	)
	cont_enteringKey.Hide()

	//-------------------------- ЭКРАН ЗАВЕРШЕНИЯ РАБОТЫ --------------------------

	img_endWork := canvas.NewImageFromFile("end_work.png")
	box_endWork := container.NewHBox(
		layout.NewSpacer(),
		container.New(
			layout.NewGridWrapLayout(fyne.NewSize(96, 96)),
			img_endWork,
		),
		layout.NewSpacer(),
	)

	text_end := canvas.NewText("Работа над файлом завершена!", color.Black)
	text_end.Alignment = fyne.TextAlignCenter
	text_end.TextSize = 14

	cont_endScene := container.NewVBox(
		box_endWork,
		text_end,
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