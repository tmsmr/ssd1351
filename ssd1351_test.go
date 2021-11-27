package ssd1351

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	rstPin = rpio.Pin(25)
	csPin  = rpio.Pin(8)
	dcPin  = rpio.Pin(24)
)

func TestMain(m *testing.M) {
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	ret := m.Run()
	if err := rpio.Close(); err != nil {
		panic(err)
	}
	os.Exit(ret)
}

func TestInit(t *testing.T) {
	oled, err := Setup(rpio.Spi0, 0, rstPin, csPin, dcPin, false)
	if err != nil {
		t.Fatal(err)
	}

	oled.Init()

	_ = oled.Shutdown()
}

func TestOutOfBounds(t *testing.T) {
	oled, err := Setup(rpio.Spi0, 0, rstPin, csPin, dcPin, false)
	if err != nil {
		t.Fatal(err)
	}

	oled.Init()

	if err := oled.DrawPixel(128, 128, 0xFFF); err == nil {
		t.Fatal("bounds check not working")
	}

	if err := oled.DrawBlock(120, 120, 25, 25, 0xFFF); err == nil {
		t.Fatal("bounds check not working")
	}

	if err := oled.DrawPixels(127, 127, 2, 2, []uint16{0xFFF, 0xFFF, 0xFFF, 0xFFF}); err == nil {
		t.Fatal("bounds check not working")
	}

	_ = oled.Shutdown()
}

func TestDrawPixel(t *testing.T) {
	oled, err := Setup(rpio.Spi0, 0, rstPin, csPin, dcPin, false)
	if err != nil {
		t.Fatal(err)
	}

	oled.Init()

	for i := 0; i < 1000; i++ {
		x := uint8(rand.Intn(oledPixelsXY - 1))
		y := uint8(rand.Intn(oledPixelsXY - 1))
		c := uint16(rand.Intn(0xFFFF))
		_ = oled.DrawPixel(x, y, c)
		time.Sleep(10 * time.Millisecond)
	}

	_ = oled.Shutdown()
}

func TestDrawBlock(t *testing.T) {
	oled, err := Setup(rpio.Spi0, 0, rstPin, csPin, dcPin, false)
	if err != nil {
		t.Fatal(err)
	}

	oled.Init()

	var i uint8
	for i = 128; i > 0; i = i - 8 {
		c := uint16(rand.Intn(0xFFFF))
		p := uint8((128 - i) / 2)
		_ = oled.DrawBlock(p, p, i, i, c)
		time.Sleep(250 * time.Millisecond)
	}

	_ = oled.Shutdown()
}

func TestDrawPixels(t *testing.T) {
	oled, err := Setup(rpio.Spi0, 0, rstPin, csPin, dcPin, false)
	if err != nil {
		t.Fatal(err)
	}

	oled.Init()

	for c := 0; c < 256; c++ {
		pixels := make([]uint16, 128*128)
		for y := 0; y < 128; y++ {
			for x := 0; x < 128; x++ {
				pixels[x+(y*128)] = RGBto16bit(uint8(c-x*2), uint8(c-y*2), uint8(c))
			}
		}
		_ = oled.DrawPixels(0, 0, 128, 128, pixels)
	}

	_ = oled.Shutdown()
}
