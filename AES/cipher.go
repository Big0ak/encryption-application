package AES

import (
	"encoding/binary"
	"errors"
	"strings"
)

func Encrypt(key, plain, ext []byte) ([]byte, error) {
	// Ключ 128 бит (16 байт)
	if len(key) != 16 {
		return nil, errors.New("the key is not supported")
	}
	sizePlainText := len(plain) // Каждый блок 1 байт => общий размер данных в байтах
	// TODO: пока максимальный размер только int32, в дальнейшем увеличится

	// Размер не больше размерности для int32 - 15 (т.е ~ 2Гб)
	if sizePlainText > 2147483632 {
		return nil, errors.New("large size of the encrypted file")
	}

	// -------------------------------AES------------------------------------------------
	// Расширения ключа для AES
	var w [nb * (nr + 1)]uint32
	expandKey(key, &w)

	// Шифрование по блокам 16 байт
	var l, r int = 0, 16
	enc := make([]byte, 0) // Зашифрованные данные
	for ;r <= sizePlainText;{
		encBlock := encryptBlock(w, plain[l:r])
		enc = append(enc, encBlock...) // TODO: Каждый раз пересоздает слайс
		l += 16; r += 16 
	}

	// В последний блок добавляется нули, если не хватает размерности
	remains := 16 - sizePlainText % 16
	if remains != 16 {
		encBlock := encryptBlock(w, append(plain[l:sizePlainText], make([]byte, remains)...))
		enc = append(enc, encBlock...) 
	}
	// ----------------------------------------------------------------------------------

	// Размерность и расширение шифруются в послденем блоке
	// Размерность занимает 8 байт, расширение 8 байт
	// Недостающие значение заполняются нулями в начале

	// Переводится в массив [8]byte (8 байт всего)
	size := make([]byte, 8)
	binary.BigEndian.PutUint64(size, uint64(sizePlainText))
	
	// Дополнение нулями в начале, чтобы размерность расширение была 8 байт
	ext = append(make([]byte, 8 - len(ext)), ext...)
	
	encBlock := encryptBlock(w, append(size, ext...))
	enc = append(enc, encBlock...)

	return enc, nil
}

func Decrypt(key, enc []byte) ([]byte, string, error) {
	// Ключ 128 бит (16 байт)
	if len(key) != 16 {
		return nil, "", errors.New("the key is not supported")
	}
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
	expandKey(key, &w)

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