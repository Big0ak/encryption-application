package spoof

import (
	"io/ioutil"
	"net/http"
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
	resp, _ := http.Get("https://random-word-api.herokuapp.com/word")
	body, _ := ioutil.ReadAll(resp.Body)

	f, _ := os.Create(t.nameFile + t.ext)
	f.Write(body)
	return t.nameFile + t.ext
}