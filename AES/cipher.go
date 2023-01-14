package AES

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Encrypt(key []byte, pathFile string) (string, error) {
	// исходный файл
	plainFile, err := os.Open(pathFile)
	if err != nil {
		return "", err
	}
	
	f, err := plainFile.Stat()
	if err != nil {
		return "", err
	}
	sizePlainText := f.Size()

	// зашифрованный файл
	encFile, err := os.Create(CipherFile)
	if err != nil {
		return "", err
	}
	
	defer func () {
		plainFile.Close()
		encFile.Close()
	} ()

	// -------------------------------AES------------------------------------------------
	// Расширения ключа для AES
	var w [nb * (nr + 1)]uint32
	// Ключ 256 бит (32 байт)
	extended_key := sha256.Sum256(key)
	expandKey(extended_key[:], &w)

	// Шифрование по блокам 16 байт
	var n int = 0
	encBlock := make([]byte, 0, 16) // Зашифрованные данные
	plainBlock := make([]byte, 16)

	n, err = plainFile.Read(plainBlock)
	for ;err != io.EOF; {
		if n == 16 {
			encBlock = encryptBlock(w, plainBlock[:])
		} else {
			// В последний блок добавляется нули, если не хватает размерности
			encBlock = encryptBlock(w, append(plainBlock[:n], make([]byte, 16 - n)...))
		}
		encFile.Write(encBlock)

		n, err = plainFile.Read(plainBlock)
	}

	// ----------------------------------------------------------------------------------

	// Размерность и расширение шифруются в послденем блоке
	// Размерность занимает 8 байт, расширение 8 байт
	// Недостающие значение заполняются нулями в начале

	// Переводится в массив [8]byte (8 байт всего)
	size := make([]byte, 8)
	binary.BigEndian.PutUint64(size, uint64(sizePlainText))
	
	// Дополнение нулями в начале, чтобы размерность расширение была 8 байт
	ext := []byte(filepath.Ext(pathFile))
	ext = append(make([]byte, 8 - len(ext)), ext...)
	
	encBlock = encryptBlock(w, append(size, ext...))
	encFile.Write(encBlock)

	return CipherFile, nil
}

func Decrypt(key, enc []byte) ([]byte, string, error) {

	// Размер не больше размерности для int32 - 15 (т.е ~ 2Гб)
	// Размер зашифрованного сообщение всегда кратен 16 
	if len(enc) > 2147483647 || len(enc) % 16 != 0 {
		return nil, "", errors.New("large encrypted file or is it corrupted")
	}

	var sizePlainText uint64
	var ext string

	// -------------------------------AES------------------------------------------------
	// Расширения ключа для AES
	var w [nb * (nr + 1)]uint32
	// Ключ 256 бит (32 байт)
	extended_key := sha256.Sum256(key)
	expandKey(extended_key[:], &w)

	// Поиск из последнего блока: размера исходных данных и расширение файла
	plainBlock := decryptBlock(w, enc[len(enc)-16:])
	sizePlainText = binary.BigEndian.Uint64(plainBlock[0:8])
	ext = strings.Replace(string(plainBlock[8:16]), "\x00" , "", -1)

	// Расшифрование по блокам 16 байт
	var l, r uint64 = 0, 16
	plain := make([]byte, 0)
	for ;r <= sizePlainText;{
		plainBlock = decryptBlock(w, enc[l:r])
		plain = append(plain, plainBlock...) // TODO: Каждый раз пересоздает слайс
		l += 16; r += 16 
	}

	// Последний блок
	remains := int(sizePlainText % 16)
	if remains != 0 {
		plainBlock := decryptBlock(w, enc[l:r])
		plain = append(plain, plainBlock[0:remains]...) 
	}
	// ----------------------------------------------------------------------------------

	return plain, ext, nil
}