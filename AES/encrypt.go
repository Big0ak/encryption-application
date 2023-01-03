package AES

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

func mixColum(s uint32) uint32{
	return (gfMul(2, s>>24&0xff) ^ gfMul(3, s>>16&0xff) ^ s>>8&0xff ^ s&0xff)<<24 |
		(s>>24&0xff ^ gfMul(2, s>>16&0xff) ^ gfMul(3, s>>8&0xff) ^ s&0xff)<<16 |
		(s>>24&0xff ^ s>>16&0xff ^ gfMul(2, s>>8&0xff) ^ gfMul(3, s&0xff))<<8 |
		(gfMul(3, s>>24&0xff) ^ s>>16&0xff ^ s>>8&0xff ^ gfMul(2, s&0xff))
}

func encryptBlock(w [nb * (nr + 1)]uint32, plaintBlock []byte) []byte {
	// разбиение на столбцы исходного блока
	s0 := uint32(plaintBlock[0])<<24 | uint32(plaintBlock[4])<<16 | uint32(plaintBlock[8])<<8 | uint32(plaintBlock[12])
	s1 := uint32(plaintBlock[1])<<24 | uint32(plaintBlock[5])<<16 | uint32(plaintBlock[9])<<8 | uint32(plaintBlock[13])
	s2 := uint32(plaintBlock[2])<<24 | uint32(plaintBlock[6])<<16 | uint32(plaintBlock[10])<<8 | uint32(plaintBlock[14])
	s3 := uint32(plaintBlock[3])<<24 | uint32(plaintBlock[7])<<16 | uint32(plaintBlock[11])<<8 | uint32(plaintBlock[15])
	
	// AddRoundKey сразу применяется ко всему столбцу
	s0 ^= w[0]
	s1 ^= w[1]
	s2 ^= w[2]
	s3 ^= w[3]

	var t0, t1, t2, t3 uint32
	for r := 1; r < nr; r++ {
		// ShiftRows, SubBytes "subWord", mixColum, AddRoundKey
		t0 = mixColum(subWord((s0>>24&0xff)<<24 | (s1>>16&0xff)<<16 | (s2>>8&0xff)<<8 | s3&0xff)) ^ w[4*r+0]
		t1 = mixColum(subWord((s1>>24&0xff)<<24 | (s2>>16&0xff)<<16 | (s3>>8&0xff)<<8 | s0&0xff)) ^ w[4*r+1]
		t2 = mixColum(subWord((s2>>24&0xff)<<24 | (s3>>16&0xff)<<16 | (s0>>8&0xff)<<8 | s1&0xff)) ^ w[4*r+2]
		t3 = mixColum(subWord((s3>>24&0xff)<<24 | (s0>>16&0xff)<<16 | (s1>>8&0xff)<<8 | s2&0xff)) ^ w[4*r+3]

		s0, s1, s2, s3 = t0, t1, t2, t3
	}

	// поледние раунды hiftRows, SubBytes, AddRoundKey
	s0 = subWord((t0>>24&0xff)<<24 | (t1>>16&0xff)<<16 | (t2>>8&0xff)<<8 | t3&0xff) ^ w[4*nr+0]
	s1 = subWord((t1>>24&0xff)<<24 | (t2>>16&0xff)<<16 | (t3>>8&0xff)<<8 | t0&0xff) ^ w[4*nr+1]
	s2 = subWord((t2>>24&0xff)<<24 | (t3>>16&0xff)<<16 | (t0>>8&0xff)<<8 | t1&0xff) ^ w[4*nr+2]
	s3 = subWord((t3>>24&0xff)<<24 | (t0>>16&0xff)<<16 | (t1>>8&0xff)<<8 | t2&0xff) ^ w[4*nr+3]

	encBlock := make([]byte, 0)
	encBlock = append(encBlock, []byte{byte(s0>>24&0xff), byte(s1>>24&0xff), byte(s2>>24&0xff), byte(s3>>24&0xff)}...)
	encBlock = append(encBlock, []byte{byte(s0>>16&0xff), byte(s1>>16&0xff), byte(s2>>16&0xff), byte(s3>>16&0xff)}...)
	encBlock = append(encBlock, []byte{byte(s0>>8&0xff), byte(s1>>8&0xff), byte(s2>>8&0xff), byte(s3>>8&0xff)}...)
	encBlock = append(encBlock, []byte{byte(s0&0xff), byte(s1&0xff), byte(s2&0xff), byte(s3&0xff)}...)
	return encBlock
}