# SSD1351
*Golang module for SSD1351-based OLED's on a Raspberry Pi*

*This module was developed for https://www.waveshare.com/wiki/1.5inch_RGB_OLED_Module, but should work with other OLED's (128x128, RGB) based on the SSD1351 as well.*

## Requirements
- A user with access to `/dev/mem` or `/dev/gpiomem` (possibly `root`)
- The `SSD1351` connected to SPI0 (`go-rpi` only supports SPI0 at the moment)

## Example usage

| RPi GPIO       | OLED |
|----------------|------|
| MOSI0 (GPIO10) | DIN  |
| SCLK0 (GPIO11) | CLK  |
| CE0 (GPIO8)    | CS   |
| GPIO24         | DC   |
| GPIO25         | RST  |

 ```go
import (
	"time"
	"github.com/tmsmr/ssd1351"
    "github.com/stianeikeland/go-rpio/v4"
)

...

rstPin := rpio.Pin(25)
csPin := rpio.Pin(8)
dcPin := rpio.Pin(24)
oled, _ := ssd1351.Setup(rpio.Spi0, 0, rstPin, csPin, dcPin, true)
oled.Init()
color := ssd1351.RGBto16bit(0xFF, 0x00, 0x00)
_ = oled.DrawBlock(0, 0, 128, 128, color)
time.Sleep(500 * time.Millisecond)
oled.Shutdown()
 ```
