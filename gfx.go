package ssd1351

func RGBto16bit(r uint8, g uint8, b uint8) uint16 {
	// 0bRRRRR-GGGGGG-BBBBB
	return uint16(r>>3)<<11 | uint16(g>>2)<<5 | uint16(b>>3)
}
