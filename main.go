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

var (
	sourceFileLoaded = false // сигнал о том, что файл загружен
	endWork = false          // сигнал о том, что шифрование/дешифрование прошло над файлом успешно => финальное окно
	newWork = false 		 // возвращение на главный экран
)

func main() {
	// ------------------- test value ------------------------------------
	// text_byte := []byte{ 0x00, 0x11, 0x22, 0x33,
	// 					0x44, 0x55, 0x66, 0x77,
	// 					0x88, 0x99, 0xaa, 0xbb,
	// 					0xcc, 0xdd, 0xee, 0xff}
	// key_byte := []byte{ 0x00, 0x01, 0x02, 0x03,
	// 					0x04, 0x05, 0x06, 0x07,
	// 					0x08, 0x09, 0x0a, 0x0b,
	// 					0x0c, 0x0d, 0x0e, 0x0f,
	// 					0x10, 0x11, 0x12, 0x13,
	// 					0x14, 0x15, 0x16, 0x17,
	// 					0x18, 0x19, 0x1a, 0x1b,
	// 					0x1c, 0x1d, 0x1e, 0x1f,};
	// enc_byte, _ := AES.Encrypt(key_byte[:], text_byte[:], []byte("fe3"))
	// plain_byte, _, _ := AES.Decrypt(key_byte[:], enc_byte[:])
	// --------------------------------------------------------------------

	var (
		key      = ""
		nameFile = ""
		ext      = ""
		pathFile = ""
	)

	a := app.New()
	w := a.NewWindow(" ")
	w.Resize(fyne.NewSize(250, 350))
	w.CenterOnScreen()   // окно по центру экрана
	w.SetFixedSize(true) // нельзя менять размер
	w.SetMaster()        // главное окно

	//-------------------------- ЭКРАН ЗАГРУЗКИ ФАЙЛА --------------------------

	text_title := canvas.NewText("Шифрование файлов AES", color.Black)
	text_title.Alignment = fyne.TextAlignCenter
	text_title.TextSize = 18
	text_title.Resize(fyne.NewSize(240, 30))
	text_title.Move(fyne.NewPos(0, 30))

	img_fileUpload := canvas.NewImageFromResource(resourceFileUploadPng)
	img_fileUpload.Resize(fyne.NewSize(80, 80))
	img_fileUpload.Move(fyne.NewPos(85, 135))

	btn_openFile := container.New(
		layout.NewMaxLayout(),

		widget.NewButton(" ", func() {
			w2 := a.NewWindow(" ")
			w2.Resize(fyne.NewSize(520, 360))
			w2.CenterOnScreen()
			w2.SetFixedSize(true)
			dialog.ShowFileOpen(
				func(uc fyne.URIReadCloser, err error) {
					nameFile = uc.URI().Name()
					ext = uc.URI().Extension()
					pathFile = uc.URI().Path()
					sourceFileLoaded = true
					w2.Close()
				},
				w2,
			)
			w2.Show()
		}),
		canvas.NewHorizontalGradient(color.RGBA{81, 81, 81, 255}, color.Transparent),
	)
	btn_openFile.Resize(fyne.NewSize(240, 100))
	btn_openFile.Move(fyne.NewPos(0, 125))

	text_footer := canvas.NewText("Загрузите файл, чтобы продолжить", color.Black)
	text_footer.Alignment = fyne.TextAlignCenter
	text_footer.TextSize = 13
	text_footer.Resize(fyne.NewSize(240, 30))
	text_footer.Move(fyne.NewPos(0, 272))

	cont_UploadFile := container.NewWithoutLayout(
		text_title,
		btn_openFile,
		img_fileUpload,
		text_footer,
	)

	//-------------------------- ЭКРАН ВВОДА КЛЮЧА --------------------------

	img_fileIcon := canvas.NewImageFromResource(resourceFileIconPng)
	box_fileIcon := container.NewHBox(
		layout.NewSpacer(),
		container.New(
			layout.NewGridWrapLayout(fyne.NewSize(96, 96)),
			img_fileIcon,
		),
		layout.NewSpacer(),
	)

	field_nameFile := binding.NewString()
	box_nameFile := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabelWithData(field_nameFile),
		layout.NewSpacer(),
	)

	field_typeWork := binding.NewString()
	box_typeWork := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabelWithData(field_typeWork),
		layout.NewSpacer(),
	)

	field_checkKey := binding.NewString()
	box_checkKey := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabelWithData(field_checkKey),
		layout.NewSpacer(),
	)

	wid_inputKey := widget.NewPasswordEntry()
	wid_inputKey.SetPlaceHolder("Введите ключ")
	// Валидация ключа
	wid_inputKey.Validator = func(input string) error {
		if len([]rune(input)) < 8 {
			return errors.New("")
		}
		return nil
	}

	pd_crypto := widget.NewProgressBarInfinite()
	pd_crypto.Stop()
	pd_crypto.Hide()

	btn_crypto := widget.NewButton("Crypto", func() {
		// Проверка валидации ключа
		if err := wid_inputKey.Validate(); err == nil {

			pd_crypto.Start()
			pd_crypto.Show()

			key = wid_inputKey.Text
			if ext == AES.ExpCrypto {
				// Расшифровка
				_, err = AES.Decrypt([]byte(key), pathFile)
			} else {
				// Шифрование
				_, err = AES.Encrypt([]byte(key), pathFile)
			}
			if err != nil {
				dialog.ShowError(err, w)
			} else {
				endWork = true
			}			
		} else {
			field_checkKey.Set("Ключ минимум 8 символов")
		}
	})

	cont_enteringKey := container.NewVBox(
		box_fileIcon,
		box_nameFile,
		box_typeWork,
		box_checkKey,
		wid_inputKey,
		btn_crypto,
		pd_crypto,
	)
	cont_enteringKey.Hide()

	//-------------------------- ЭКРАН ЗАВЕРШЕНИЯ РАБОТЫ --------------------------

	img_endWork := canvas.NewImageFromResource(resourceEndWorkPng)
	img_endWork.Resize(fyne.NewSize(80,80))
	img_endWork.Move(fyne.NewPos(80,30))

	text_end := canvas.NewText("Работа над файлом завершена!", color.Black)
	text_end.Alignment = fyne.TextAlignCenter
	text_end.TextSize = 13
	text_end.Resize(fyne.NewSize(240, 30))
	text_end.Move(fyne.NewPos(0, 120))

	btn_newWork := container.NewVBox(
		layout.NewSpacer(),
		widget.NewButton("На главную", func() {
			newWork = true
		}),
	)
	btn_newWork.Resize(fyne.NewSize(160,20))
	btn_newWork.Move(fyne.NewPos(45,310))

	cont_endScene := container.NewWithoutLayout(
		img_endWork,
		text_end,
		btn_newWork,
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
				nameAll := []rune(nameFile)
				if len(nameAll) > 30 {
					field_nameFile.Set(string(nameAll[0:27]) + "...")
				} else {
					field_nameFile.Set(nameFile)
				}
				if ext == AES.ExpCrypto {
					field_typeWork.Set("Введите ключ для расшифровки")
				} else {
					field_typeWork.Set("Введите ключ для шифрования")
				}
				cont_enteringKey.Show()
				sourceFileLoaded = false
			}

			if endWork {
				cont_enteringKey.Hide()
				cont_endScene.Show()
				pd_crypto.Stop()
				pd_crypto.Hide()
				endWork = false
			}

			if newWork {
				cont_endScene.Hide()
				cont_UploadFile.Show()
				wid_inputKey.SetText("")
				field_checkKey.Set("")
				newWork = false
			}
		}
	}()

	a.Run()
}
