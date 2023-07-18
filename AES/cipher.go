package AES

import (
	"crypto/sha256"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Big0ak/AES/spoof"
)

func Encrypt(key []byte, pathFile string) (string, error) {
	// Исходный файл
	plainFile, err := os.Open(pathFile)
	if err != nil {
		return "", err
	}
	
	f, err := plainFile.Stat()
	if err != nil {
		return "", err
	}
	// Размер исходного файла в байтах
	sizePlainFile := f.Size()

	// Создание зашифрованного файла
	encFile, err := os.Create(CipherFile)
	if err != nil {
		return "", err
	}

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
	binary.BigEndian.PutUint64(size, uint64(sizePlainFile))
	
	// Дополнение нулями в начале, чтобы размерность расширение была 8 байт
	ext := []byte(filepath.Ext(pathFile))
	ext = append(make([]byte, 8 - len(ext)), ext...)
	
	encBlock = encryptBlock(w, append(size, ext...))
	encFile.Write(encBlock)

	plainFile.Close()
	encFile.Close()

	return CipherFile, nil
}

func Decrypt(key []byte, pathFile string) (string, error) {
	// Исходный файл
	encFile, err := os.Open(pathFile)
	if err != nil {
		return "", err
	}

	f, err := encFile.Stat()
	if err != nil {
		return "", err
	}
	// Размер исходного файла в байтах
	sizeEncFile := f.Size()

	// в программе заложено, если файл не кратен 16 блокам => он поврежден (см. Encrypt)
	if sizeEncFile < 32 || sizeEncFile % 16 != 0 {
		s := spoof.NewSpoof(DecryptFile)
		return s.GenerateSpoofFile(), nil
		//return "", errors.New("Зашифрованный файл поврежден")
	}
	
	// -------------------------------AES------------------------------------------------
	// Расширения ключа для AES
	var w [nb * (nr + 1)]uint32
	// Ключ 256 бит (32 байт)
	extended_key := sha256.Sum256(key)
	expandKey(extended_key[:], &w)

	// Поиск из последнего блока: размера исходных данных и расширение файла
	encBlock := make([]byte, 16)
	plainBlock := make([]byte, 0, 16)

	encFile.ReadAt(encBlock, sizeEncFile-16)

	plainBlock = decryptBlock(w, encBlock)
	sizePlainFile := binary.BigEndian.Uint64(plainBlock[0:8])
	ext := strings.Replace(string(plainBlock[8:16]), "\x00" , "", -1)

	// если размерности файлов не совпадают после расшифровки => неверный ключ
	// в открытом тексте может быть меньше байт чем в зашифрованном
	// потому что в зашифрованном добавляется дополнительный блок до ровных 16 байт
	diff := (sizeEncFile - 16) - int64(sizePlainFile)
	if diff >= 16 || diff < 0 {
		s := spoof.NewSpoof(DecryptFile)
		return s.GenerateSpoofFile(), nil
		//return "", errors.New("Неверный ключ")
	}

	plainFile, err := os.Create(DecryptFile + ext)
	if err != nil {
		return "", err
	}

	// Расшифрование по блокам 16 байт
	var offset uint64 = 0
	for ;offset + 16 <= sizePlainFile;{
		encFile.ReadAt(encBlock, int64(offset)) 
		plainBlock = decryptBlock(w, encBlock)
		plainFile.Write(plainBlock)
		offset += 16 
	}

	// Последний блок
	remains := int(sizePlainFile % 16)
	if remains != 0 {
		encFile.ReadAt(encBlock, int64(offset))
		plainBlock := decryptBlock(w, encBlock)
		plainFile.Write(plainBlock[0:remains]) 
	}
	// ----------------------------------------------------------------------------------

	plainFile.Close()
	encFile.Close()

	return DecryptFile + ext, nil
}