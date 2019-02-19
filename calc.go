package dht11

func calcT(i, f byte) (res float64) {
	res = float64(i)
	res += float64(f&0x0f) * 0.1
	return
}

func calcH(i, f byte) (res float64) {
	res = float64(i)
	res += float64(f) * 0.1
	return
}
