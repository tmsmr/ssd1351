package ssd1351

import (
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

const (
	spiSpeedHz    = 9000000 // 9 MHz
	spiModePol    = 0       // SPI mode 0 (POL)
	spiModePha    = 0       // SPI mode 0 (PHA)
	spiCSPolarity = 0       // CS active low

	ssd1351CmdCommandLock    = 0xFD // SSD1351.pdf: 10.1.23 Set Command Lock (FDh)
	ssd1351CmdDisplayOff     = 0xAE // SSD1351.pdf: 10.1.10 Set Sleep mode ON/OFF (AEh / AFh)
	ssd1351CmdDisplayAllOff  = 0xA4 // SSD1351.pdf: 10.1.8 Set Display Mode (A4h ~ A7h)
	ssd1351CmdSetColumn      = 0x15 // SSD1351.pdf: 10.1.1 Set Column Address (15h)
	ssd1351CmdSetRow         = 0x75 // SSD1351.pdf: 10.1.2 Set Row Address (75h)
	ssd1351CmdClockDiv       = 0xB3 // SSD1351.pdf: 10.1.13 Set Front Clock Divider / Oscillator Frequency (B3h)
	ssd1351CmdMuxRatio       = 0xCA // SSD1351.pdf: 10.1.22 Set Multiplex Ratio (CAh)
	ssd1351CmdSetRemap       = 0xA0 // SSD1351.pdf: 10.1.5 Set Re-map & Dual COM Line Mode (A0h)
	ssd1351CmdStartLine      = 0xA1 // SSD1351.pdf: 10.1.6 Set Display Start Line (A1h)
	ssd1351CmdDisplayOffset  = 0xA2 // SSD1351.pdf: 10.1.7 Set Display Offset (A2h)
	ssd1351CmdFunctionSelect = 0xAB // SSD1351.pdf: 10.1.9 Set Function selection (ABh)
	ssd1351CmdSetVSL         = 0xB4 // SSD1351.pdf: Page 34, B4
	ssd1351CmdContrastABC    = 0xC1 // SSD1351.pdf: 10.1.20 Set Contrast Current for Color A,B,C (C1h)
	ssd1351CmdContrastMaster = 0xC7 // SSD1351.pdf: 10.1.21 Master Contrast Current Control (C7h)
	ssd1351CmdPrecharge      = 0xB1 // SSD1351.pdf: 10.1.11 Set Phase Length (B1h)
	ssd1351CmdDisplayEnhance = 0xB2 // SSD1351.pdf: 10.1.12 Display Enhancement (B2h)
	ssd1351CmdPrecharge2     = 0xB6 // SSD1351.pdf: 10.1.15 Set Second Pre-charge period (B6h)
	ssd1351CmdPrechargeLevel = 0xBB // SSD1351.pdf: 10.1.18 Set Pre-charge voltage (BBh)
	ssd1351CmdVCOMH          = 0xBE // SSD1351.pdf: 10.1.19 Set VCOMH Voltage (BEh)
	ssd1351CmdNormalDisplay  = 0xA6 // SSD1351.pdf: 10.1.8 Set Display Mode (A4h ~ A7h)
	ssd1351CmdDisplayOn      = 0xAF // SSD1351.pdf: 10.1.10 Set Sleep mode ON/OFF (AEh / AFh)
	ssd1351CmdWriteRAM       = 0x5C // SSD1351.pdf: 10.1.3 Write RAM Command (5Ch)

	oledPixelsXY = 128
)

// Setup opens the connection to the OLED using four-wire SPI
// dev: The rpio.SpiDev to use
// slave: The slave chip number to use
// rstPin: Reset rpio.Pin, needed for initialization
// csPin: Chip select rpio.Pin, needed to switch between command/data (8.1.3 MCU Serial Interface (4-wire SPI), 8.4 Command Decoder)
// dcPin: Data/command select rpio.Pin, needed to switch between command/data (8.1.3 MCU Serial Interface (4-wire SPI), 8.4 Command Decoder)
// openGpio: Shall this package call rpio.Open()/rpio.Close() or is this performed outside?
func Setup(dev rpio.SpiDev, slave uint8, rstPin rpio.Pin, csPin rpio.Pin, dcPin rpio.Pin, openGpio bool) (*SSD1351, error) {
	if openGpio {
		if err := rpio.Open(); err != nil {
			return nil, err
		}
	}
	rstPin.Output()
	csPin.Output()
	dcPin.Output()
	if err := rpio.SpiBegin(dev); err != nil {
		return nil, err
	}
	rpio.SpiChipSelect(slave)
	rpio.SpiChipSelectPolarity(slave, spiCSPolarity)
	rpio.SpiSpeed(spiSpeedHz)
	rpio.SpiMode(spiModePol, spiModePha)
	return &SSD1351{dev: dev, rstPin: rstPin, csPin: csPin, dcPin: dcPin, openGpio: openGpio}, nil
}

type CMD struct {
	cmd  uint8
	data []uint8
}

// defConfSeq returns the the default configuration sequence for the SSD1351
// This sequence is used in Waveshare's example for Python (https://www.waveshare.com/wiki/1.5inch_OLED_Module)
// I don't comprehend anything at the moment and may review this someday...
func defConfSeq() []CMD {
	return []CMD{
		{ssd1351CmdCommandLock, []uint8{0x12}},                 // reset mcu protection status
		{ssd1351CmdCommandLock, []uint8{0xB1}},                 // make commands A2,B1,B3,BB,BE,C1 accessible
		{ssd1351CmdDisplayOff, nil},                            // instruct display to sleep
		{ssd1351CmdDisplayAllOff, nil},                         // turn the display off
		{ssd1351CmdSetColumn, []uint8{0x00, oledPixelsXY - 1}}, // why ???
		{ssd1351CmdSetRow, []uint8{0x00, oledPixelsXY - 1}},    // why ???
		{ssd1351CmdClockDiv, []uint8{0b11110001}},              // set oscillator freq. to 0b1111, divide by 1
		{ssd1351CmdMuxRatio, []uint8{0x7F}},                    // set mux ratio to 127
		{ssd1351CmdSetRemap, []uint8{0b01110100}},              // scan from COM[N-1] to COM0, color sequence C -> B -> A
		{ssd1351CmdStartLine, []uint8{0x00}},                   // set start line to 0
		{ssd1351CmdDisplayOffset, []uint8{0x00}},               // set offset to 0
		{ssd1351CmdFunctionSelect, []uint8{0x01}},              // internal VDD regulator, SPI interface
		{ssd1351CmdSetVSL, []uint8{0b10100000, 0xB5, 0x55}},    // use external VSL
		{ssd1351CmdContrastABC, []uint8{0xC8, 0x80, 0xC0}},     // set contrast for A, B, C
		{ssd1351CmdContrastMaster, []uint8{0x0F}},              // max contrast factor
		{ssd1351CmdPrecharge, []uint8{0b00110010}},             // phase 1: 5 DCLKs, phase 2: 3 DCLKs
		{ssd1351CmdDisplayEnhance, []uint8{0xA4, 0x00, 0x00}},  // enhance display performance
		{ssd1351CmdPrechargeLevel, []uint8{0x17}},              // set pre-charge voltage to ~70%
		{ssd1351CmdPrecharge2, []uint8{0x01}},                  // second pre-charge period: 1 DCLK
		{ssd1351CmdVCOMH, []uint8{0x05}},                       // set VCOMMH to 0.82 * VCC
		{ssd1351CmdNormalDisplay, nil},                         // set display to "normal" operation
	}
}

type SSD1351 struct {
	dev      rpio.SpiDev
	rstPin   rpio.Pin
	csPin    rpio.Pin
	dcPin    rpio.Pin
	openGpio bool
}

func (s *SSD1351) Init() {
	// reset SSD1351
	s.csPin.Low()
	s.rstPin.Low()
	time.Sleep(1 * time.Millisecond)
	s.rstPin.High()
	time.Sleep(300 * time.Millisecond)
	// configure SSD1351
	for _, cmd := range defConfSeq() {
		s.send(cmd)
	}
	// clear display and activate
	s.ClearScreen()
	s.txCmd(ssd1351CmdDisplayOn)
}

func (s *SSD1351) tx(data ...byte) {
	rpio.SpiTransmit(data...)
}

func (s *SSD1351) txCmd(cmd byte) {
	s.csPin.Low()
	s.dcPin.Low()
	s.tx(cmd)
	s.csPin.High()
}

func (s *SSD1351) txData(data ...byte) {
	s.csPin.Low()
	s.dcPin.High()
	s.tx(data...)
	s.csPin.High()
}

func (s *SSD1351) send(cmd CMD) {
	s.txCmd(cmd.cmd)
	s.txData(cmd.data...)
}

func (s *SSD1351) setGDDRAMAddr(c1 uint8, c2 uint8, r1 uint8, r2 uint8) {
	s.send(CMD{
		cmd:  ssd1351CmdSetColumn,
		data: []uint8{c1, c2},
	})
	s.send(CMD{
		cmd:  ssd1351CmdSetRow,
		data: []uint8{r1, r2},
	})
	s.txCmd(ssd1351CmdWriteRAM)
}

func (s *SSD1351) ClearScreen() {
	s.setGDDRAMAddr(0, oledPixelsXY-1, 0, oledPixelsXY-1)
	clearBytes := make([]uint8, oledPixelsXY*oledPixelsXY*2)
	for i := range clearBytes {
		clearBytes[i] = 0x00
	}
	s.txData(clearBytes...)
}

func (s *SSD1351) DrawPixel(x uint8, y uint8, color uint16) {
	s.setGDDRAMAddr(x, x, y, y)
	s.txData([]uint8{uint8(color >> 8), uint8(color & 0xFF)}...)
}

func (s *SSD1351) Shutdown() error {
	s.txCmd(ssd1351CmdDisplayOff)
	time.Sleep(100 * time.Millisecond)
	rpio.SpiEnd(s.dev)
	if s.openGpio {
		return rpio.Close()
	}
	return nil
}
