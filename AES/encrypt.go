package AES

import (
	"encoding/binary"
)

func gfMul(a, b uint32) uint32 {
	i := a
	j := b
	res := uint32(0)
	for k := uint32(1); k < 0x100 && j != 0; k <<= 1 {
		//k == 1<<n, i == b * xⁿ
		if j&k != 0 {
			res ^= i
			j ^= k
		}

		i <<= 1
		if i&0x100 != 0 {
			i ^= 1<<8 | 1<<4 | 1<<3 | 1<<1 | 1<<0 // x⁸ + x⁴ + x³ + x + 1
		}
	}
	return res
}

//------------------------------------------------------------------------------------------------------
//------------------------------------------ENCRYPT-----------------------------------------------------
//------------------------------------------------------------------------------------------------------

func mixColum(s uint32) uint32 {
	return (gfMul(2, s>>24&0xff)^gfMul(3, s>>16&0xff)^s>>8&0xff^s&0xff)<<24 |
		(s>>24&0xff^gfMul(2, s>>16&0xff)^gfMul(3, s>>8&0xff)^s&0xff)<<16 |
		(s>>24&0xff^s>>16&0xff^gfMul(2, s>>8&0xff)^gfMul(3, s&0xff))<<8 |
		(gfMul(3, s>>24&0xff)^s>>16&0xff^s>>8&0xff^gfMul(2, s&0xff))
}

func subBytesCol(s uint32) uint32 {
	return uint32(sbox[s>>28][s>>24&0xf])<<24 |
		uint32(sbox[s>>20&0xf][s>>16&0xf])<<16 |
		uint32(sbox[s>>12&0xf][s>>8&0xf])<<8 |
		uint32(sbox[s>>4&0xf][s&0xf])
}

func encryptBlock(w [nb * (nr + 1)]uint32, plainBlock []byte) []byte {
	// разбиение на столбцы исходного блока (как в официальной документации)
	s0 := binary.BigEndian.Uint32(plainBlock[0:4])
	s1 := binary.BigEndian.Uint32(plainBlock[4:8])
	s2 := binary.BigEndian.Uint32(plainBlock[8:12])
	s3 := binary.BigEndian.Uint32(plainBlock[12:16])

	// AddRoundKey сразу применяется ко всему столбцу
	s0 ^= w[0]
	s1 ^= w[1]
	s2 ^= w[2]
	s3 ^= w[3]

	var t0, t1, t2, t3 uint32
	for r := 1; r < nr; r++ {
		// все преобразования со столбцами
		// ShiftRows, SubBytes "subBytesCol", mixColum, AddRoundKey
		t0 = mixColum(subBytesCol((s0>>24&0xff)<<24|(s1>>16&0xff)<<16|(s2>>8&0xff)<<8|s3&0xff)) ^ w[4*r]
		t1 = mixColum(subBytesCol((s1>>24&0xff)<<24|(s2>>16&0xff)<<16|(s3>>8&0xff)<<8|s0&0xff)) ^ w[4*r+1]
		t2 = mixColum(subBytesCol((s2>>24&0xff)<<24|(s3>>16&0xff)<<16|(s0>>8&0xff)<<8|s1&0xff)) ^ w[4*r+2]
		t3 = mixColum(subBytesCol((s3>>24&0xff)<<24|(s0>>16&0xff)<<16|(s1>>8&0xff)<<8|s2&0xff)) ^ w[4*r+3]

		s0, s1, s2, s3 = t0, t1, t2, t3
	}

	// поcледние раунды ShiftRows, SubBytes, AddRoundKey
	s0 = subBytesCol((t0>>24&0xff)<<24|(t1>>16&0xff)<<16|(t2>>8&0xff)<<8|t3&0xff) ^ w[4*nr]
	s1 = subBytesCol((t1>>24&0xff)<<24|(t2>>16&0xff)<<16|(t3>>8&0xff)<<8|t0&0xff) ^ w[4*nr+1]
	s2 = subBytesCol((t2>>24&0xff)<<24|(t3>>16&0xff)<<16|(t0>>8&0xff)<<8|t1&0xff) ^ w[4*nr+2]
	s3 = subBytesCol((t3>>24&0xff)<<24|(t0>>16&0xff)<<16|(t1>>8&0xff)<<8|t2&0xff) ^ w[4*nr+3]

	encBlock := make([]byte, len(plainBlock))
	// сборка из столбцов в зашифрованный блок
	binary.BigEndian.PutUint32(encBlock[0:4], s0)
	binary.BigEndian.PutUint32(encBlock[4:8], s1)
	binary.BigEndian.PutUint32(encBlock[8:12], s2)
	binary.BigEndian.PutUint32(encBlock[12:16], s3)
	return encBlock
}

//------------------------------------------------------------------------------------------------------
//------------------------------------------DECRYPT-----------------------------------------------------
//------------------------------------------------------------------------------------------------------

func invMixColum(s uint32) uint32 {
	return (gfMul(0xe, s>>24&0xff)^gfMul(0xb, s>>16&0xff)^gfMul(0xd, s>>8&0xff)^gfMul(0x9, s&0xff))<<24 |
		(gfMul(0x9, s>>24&0xff)^gfMul(0xe, s>>16&0xff)^gfMul(0xb, s>>8&0xff)^gfMul(0xd, s&0xff))<<16 |
		(gfMul(0xd, s>>24&0xff)^gfMul(0x9, s>>16&0xff)^gfMul(0xe, s>>8&0xff)^gfMul(0xb, s&0xff))<<8 |
		(gfMul(0xb, s>>24&0xff)^gfMul(0xd, s>>16&0xff)^gfMul(0x9, s>>8&0xff)^gfMul(0xe, s&0xff))
}

func invSubBytesCol(s uint32) uint32 {
	return uint32(invSbox[s>>28][s>>24&0xf])<<24 |
		uint32(invSbox[s>>20&0xf][s>>16&0xf])<<16 |
		uint32(invSbox[s>>12&0xf][s>>8&0xf])<<8 |
		uint32(invSbox[s>>4&0xf][s&0xf])
}

func decryptBlock(w [nb * (nr + 1)]uint32, encBlock []byte) []byte {
	// разбиение на столбцы исходного блока (как в официальной документации)
	s0 := binary.BigEndian.Uint32(encBlock[0:4])
	s1 := binary.BigEndian.Uint32(encBlock[4:8])
	s2 := binary.BigEndian.Uint32(encBlock[8:12])
	s3 := binary.BigEndian.Uint32(encBlock[12:16])

	// AddRoundKey сразу применяется ко всему столбцу
	s0 ^= w[4*nr]
	s1 ^= w[4*nr+1]
	s2 ^= w[4*nr+2]
	s3 ^= w[4*nr+3]

	var t0, t1, t2, t3 uint32
	for r := nr - 1; r > 0; r-- {
		// все преобразования со столбцами
		// InvShiftBytes, InvSubBytes, AddRoundKey, InvMixColum
		t0 = invMixColum(invSubBytesCol((s0>>24&0xff)<<24|(s3>>16&0xff)<<16|(s2>>8&0xff)<<8|s1&0xff) ^ w[4*r])
		t1 = invMixColum(invSubBytesCol((s1>>24&0xff)<<24|(s0>>16&0xff)<<16|(s3>>8&0xff)<<8|s2&0xff) ^ w[4*r+1])
		t2 = invMixColum(invSubBytesCol((s2>>24&0xff)<<24|(s1>>16&0xff)<<16|(s0>>8&0xff)<<8|s3&0xff) ^ w[4*r+2])
		t3 = invMixColum(invSubBytesCol((s3>>24&0xff)<<24|(s2>>16&0xff)<<16|(s1>>8&0xff)<<8|s0&0xff) ^ w[4*r+3])

		s0, s1, s2, s3 = t0, t1, t2, t3
	}

	// поcледние раунды InvShiftRows, InvSubBytes, AddRoundKey
	s0 = invSubBytesCol((t0>>24&0xff)<<24|(t3>>16&0xff)<<16|(t2>>8&0xff)<<8|t1&0xff) ^ w[0]
	s1 = invSubBytesCol((t1>>24&0xff)<<24|(t0>>16&0xff)<<16|(t3>>8&0xff)<<8|t2&0xff) ^ w[1]
	s2 = invSubBytesCol((t2>>24&0xff)<<24|(t1>>16&0xff)<<16|(t0>>8&0xff)<<8|t3&0xff) ^ w[2]
	s3 = invSubBytesCol((t3>>24&0xff)<<24|(t2>>16&0xff)<<16|(t1>>8&0xff)<<8|t0&0xff) ^ w[3]

	plainBlock := make([]byte, len(encBlock))
	// сборка из столбцов в расшифрованный блок
	binary.BigEndian.PutUint32(plainBlock[0:4], s0)
	binary.BigEndian.PutUint32(plainBlock[4:8], s1)
	binary.BigEndian.PutUint32(plainBlock[8:12], s2)
	binary.BigEndian.PutUint32(plainBlock[12:16], s3)
	return plainBlock
}