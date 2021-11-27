package ssd1351

// RGBto16bit transforms 24-bit (R, G, B) colors to 16-bit colors
func RGBto16bit(r uint8, g uint8, b uint8) uint16 {
	// 0bRRRRR-GGGGGG-BBBBB
	return uint16(r>>3)<<11 | uint16(g>>2)<<5 | uint16(b>>3)
}
