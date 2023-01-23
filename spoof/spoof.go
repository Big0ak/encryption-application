package spoof

import (
	"math/rand"
	"time"
	"encoding/json"
	"net/http"
)

type ISpoofFile interface {
	GenerateSpoofFile() string
}

type Spoof struct {
	ISpoofFile
}

// файл подмены с заданным расширением
func NewSpoofExt(ext, nameFile string) Spoof {
	switch ext {

	case ".jpeg":
		return Spoof{
			ISpoofFile: newFileImage(ext, nameFile),
		}
	case ".docx":
		return Spoof{
			ISpoofFile: newFileDOCX(ext, nameFile),
		}
	default:
		return Spoof{
			ISpoofFile: newFileTXT(".txt", nameFile),
		}
	}
}

// расширения файла подмены выбирается рандомно из существующих
func NewSpoof(nameFile string) Spoof {
	rand.Seed(time.Now().UnixNano())
	ext := extBase[rand.Intn(len(extBase))]

	return NewSpoofExt(ext, nameFile)
}

// запрос по api рандомного текста
func requestText() (string, error) {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url_activity, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	if text, ok := result["activity"]; ok {
		return text.(string), nil
	}
	return "", err
}
