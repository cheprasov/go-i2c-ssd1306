//package i2c_ssd1306
package main

import (
    "github.com/cheprasov/go-i2c"
    "log"
    "time"
)

const (
    SSD1306_COMMAND    = 0x80
    SSD1306_DATA       = 0xC0
    SSD1306_MULTI_DATA = 0x40
)
const (
    SSD1306_I2C_ADDRESS           = 0x3C // 011110+SA0+RW - 0x3C or 0x3D
    SSD1306_SET_CONTRAST          = 0x81
    SSD1306_DISPLAY_ALL_ON_RESUME = 0xA4
    SSD1306_DISPLAYALLON          = 0xA5
    SSD1306_NORMAL_DISPLAY        = 0xA6
    SSD1306_INVERTDISPLAY         = 0xA7
    SSD1306_DISPLAY_OFF           = 0xAE
    SSD1306_DISPLAY_ON            = 0xAF
    SSD1306_SET_DISPLAY_OFFSET    = 0xD3
    SSD1306_SET_COM_PINS          = 0xDA
    SSD1306_SET_VCOM_DETECT       = 0xDB
    SSD1306_SET_DISPLAY_CLOCK_DIV = 0xD5
    SSD1306_SET_PRECHARGE         = 0xD9
    SSD1306_SET_MULTIPLEX         = 0xA8
    SSD1306_SETLOWCOLUMN          = 0x00
    SSD1306_SETHIGHCOLUMN         = 0x10
    SSD1306_SET_START_LINE        = 0x40
    SSD1306_MEMORY_MODE           = 0x20
    SSD1306_COLUMNADDR            = 0x21
    SSD1306_PAGEADDR              = 0x22
    SSD1306_COMSCANINC            = 0xC0
    SSD1306_COM_SCAN_DEC          = 0xC8
    SSD1306_SEG_REMAP             = 0xA0
    SSD1306_CHARGE_PUMP           = 0x8D
    SSD1306_EXTERNALVCC           = 0x1
    SSD1306_SWITCHCAPVCC          = 0x2

    // Scrolling constants
    SSD1306_ACTIVATE_SCROLL                      = 0x2F
    SSD1306_DEACTIVATE_SCROLL                    = 0x2E
    SSD1306_SET_VERTICAL_SCROLL_AREA             = 0xA3
    SSD1306_RIGHT_HORIZONTAL_SCROLL              = 0x26
    SSD1306_LEFT_HORIZONTAL_SCROLL               = 0x27
    SSD1306_VERTICAL_AND_RIGHT_HORIZONTAL_SCROLL = 0x29
    SSD1306_VERTICAL_AND_LEFT_HORIZONTAL_SCROLL  = 0x2A
)

type SSD1306 struct {
    i2cAddress uint8
    i2cBus     uint8
    i2cConnect *i2c.I2C
    width      uint8
    height     uint8
    buffer     [][]bool
}

// Create and init new instance of SSD1306 Oled Display
func NewSSD1306(i2cAddress, i2cBus, width, height uint8) (*SSD1306, error) {
    var err error

    ssd1306 := &SSD1306{
        i2cAddress: i2cAddress,
        i2cBus:     i2cBus,
        width:      width,
        height:     height,
        buffer:     [][]bool{},
    }
    err = ssd1306.initConnection()
    if err != nil {
        return nil, err
    }

    err = ssd1306.initDisplay()
    if err != nil {
        return nil, err
    }

    return ssd1306, nil
}

// Init i2c connection
func (ssd1306 *SSD1306) initConnection() error {
    connect, err := i2c.NewI2C(ssd1306.i2cAddress, int(ssd1306.i2cBus))
    if err != nil {
        return err
    }
    ssd1306.i2cConnect = connect
    return nil
}

func (ssd1306 *SSD1306) writeCommand(command int) (int, error) {
    return ssd1306.i2cConnect.WriteBytes([]byte{SSD1306_COMMAND, byte(command)})
}

func (ssd1306 *SSD1306) writeData(command int) (int, error) {
    return ssd1306.i2cConnect.WriteBytes([]byte{SSD1306_DATA, byte(command)})
}

// Init i2c Oled Display
func (ssd1306 *SSD1306) initDisplay() error {
    var err error
    _, err = ssd1306.writeCommand(SSD1306_DISPLAY_OFF)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_DISPLAY_CLOCK_DIV)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0x80)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_MULTIPLEX)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0x3F)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_DISPLAY_OFFSET)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0x0)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_START_LINE | 0x0)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_CHARGE_PUMP)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0x14)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_MEMORY_MODE)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0x00)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SEG_REMAP | 0x1)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_COM_SCAN_DEC)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_COM_PINS)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0x12)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_CONTRAST)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0xCF)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_PRECHARGE)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0xF1)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_SET_VCOM_DETECT)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(0x40)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_DISPLAY_ALL_ON_RESUME)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_NORMAL_DISPLAY)
    if err != nil {
        return err
    }
    _, err = ssd1306.writeCommand(SSD1306_DISPLAY_ON);
    if err != nil {
        return err
    }

    return nil
}

func main() {
    oled, err := NewSSD1306(0x3C, 1, 128, 64)
    if err != nil {
        log.Fatal(err)
    }

    oled.writeCommand(0x20)
    oled.writeCommand(2)
    oled.writeCommand(0xB2)
    oled.writeCommand(0x00)
    oled.writeCommand(0x12)
    oled.writeData(0x7)
    oled.writeData(0x7)
    oled.writeData(0x7)
    oled.writeData(0x7)

    time.Sleep(10 * time.Second)

    oled.writeCommand(SSD1306_DISPLAY_OFF)
}