package spoof

import (
	"io"
	"net/http"
	"os"
)

type FileImage struct {
	nameFile string
	ext      string
}

func newFileImage(ext, nameFile string) *FileImage {
	return &FileImage{
		ext:      ext,
		nameFile: nameFile,
	}
}

func (im *FileImage) GenerateSpoofFile() string {
	f, _ := os.Create(im.nameFile + im.ext)
	
	image, err := im.requestImage()
	if err != nil {
		f.Write(imageDefold)
	} else {
		_, _ = io.Copy(f, image)
	}

	return im.nameFile + im.ext
}

func (im *FileImage) requestImage() (io.ReadCloser, error) {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url_randomFaces, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}