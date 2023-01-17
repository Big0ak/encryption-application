package spoof

import (
	"github.com/srdolor/docx"
)

type FileDOCX struct {
	nameFile string
	ext      string
}

func newFileDOCX(ext, nameFile string) *FileDOCX {
	return &FileDOCX{
		ext:      ext,
		nameFile: nameFile,
	}
}

func (doc *FileDOCX) GenerateSpoofFile() string {
	f := docx.NewFile()

	text, err := requestText()
	if err != nil {
		text = textDefold
	}
	para := f.AddParagraph()
	para.AddText(text)
	f.Save(doc.nameFile + doc.ext)

	return doc.nameFile + doc.ext
}