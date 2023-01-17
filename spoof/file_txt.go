package spoof

import (
	"os"
)

type FileTXT struct {
	nameFile string
	ext string
}

func newFileTXT (ext, nameFile string) *FileTXT {
	return &FileTXT{
		ext: ext,
		nameFile: nameFile,
	}
}

func (t *FileTXT) GenerateSpoofFile() string{
	f, _ := os.Create(t.nameFile + t.ext)

	text, err := requestText()
	if err != nil{
		text = textDefold
	}
	f.Write([]byte(text))
	return t.nameFile + t.ext
}