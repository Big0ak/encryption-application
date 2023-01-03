package AES

func Encrypt (key, plain []byte) (enc []byte){
	var w [nb*(nr+1)]uint32
	expandKey(key, &w)

	enc = encryptBlock(w, plain)

	return enc
}
