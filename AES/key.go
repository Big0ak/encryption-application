package AES

func subWord(word uint32) uint32 {
	return uint32(sbox[word>>28][word>>24&0xf])<<24 |
		 uint32(sbox[word>>20&0xf][word>>16&0xf])<<16 |
		 uint32(sbox[word>>12&0xf][word>>8&0xf])<<8 |
		 uint32(sbox[word>>4&0xf][word&0xf])
}

func rotWord(word uint32) uint32 {
	return word<<8 | word>>24
}

func expandKey(key []byte, w *[nb*(nr+1)]uint32) {
	var i int
	for i = 0; i < nk; i++ {
		w[i] = uint32(key[4*i])<<24 | uint32(key[4*i+1])<<16 | uint32(key[4*i+2])<<8 | uint32(key[4*i+3])
	}

	for ; i < nb*(nr+1); i++ {
		tmp := w[i-1]
		if i % nk == 0{
			tmp = subWord(rotWord(tmp)) ^ rcon[i/nk-1]
		} else if nk > 6 && i%nk == 4 {
			tmp = subWord(tmp)
		}
		w[i] = w[i-nk] ^ tmp
	}
}