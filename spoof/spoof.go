package spoof

import (
	"math/rand"
	"time"
)

type ISpoofFile interface {
	GenerateSpoofFile() string
}

type Spoof struct {
	ISpoofFile
}

func NewSpoofExt(ext, nameFile string) Spoof {
	switch ext {

	case ".txt":
		return Spoof{
			ISpoofFile: newFileTXT(ext, nameFile),
		}

	default:
		return Spoof{}
	}
}

func NewSpoof(nameFile string) Spoof {
	rand.Seed(time.Now().UnixNano())
	ext := extBase[rand.Intn(len(extBase))]

	return NewSpoofExt(ext, nameFile)
}
