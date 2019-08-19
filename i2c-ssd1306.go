//package i2c_oled
package main

import (
    "./font"
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

const (
    SSD1306_PAGE_START_ADDRESS_0 = 0xB0
    SSD1306_PAGE_START_ADDRESS_1 = 0xB1
    SSD1306_PAGE_START_ADDRESS_2 = 0xB2
    SSD1306_PAGE_START_ADDRESS_3 = 0xB3
    SSD1306_PAGE_START_ADDRESS_4 = 0xB4
    SSD1306_PAGE_START_ADDRESS_5 = 0xB5
    SSD1306_PAGE_START_ADDRESS_6 = 0xB6
    SSD1306_PAGE_START_ADDRESS_7 = 0xB7
)

type PixelType uint8;

const (
    PIXEL_TYPE_NULL  PixelType = 0
    PIXEL_TYPE_WHITE PixelType = 1
    PIXEL_TYPE_BLACK PixelType = 1
)

type PageAddressType struct {
    PageStart        uint8
    LowerStartColumn uint8
    UpperStartColumn uint8
}

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

    oled := &SSD1306{
        i2cAddress: i2cAddress,
        i2cBus:     i2cBus,
        width:      width,
        height:     height,
        buffer:     [][]bool{},
    }
    err = oled.initConnection()
    if err != nil {
        return nil, err
    }

    err = oled.initDisplay()
    if err != nil {
        return nil, err
    }

    return oled, nil
}

// Init i2c connection
func (oled *SSD1306) initConnection() error {
    connect, err := i2c.NewI2C(oled.i2cAddress, int(oled.i2cBus))
    if err != nil {
        return err
    }
    oled.i2cConnect = connect
    return nil
}

func (oled *SSD1306) writeCommand(command byte) (int, error) {
    return oled.i2cConnect.WriteBytes([]byte{SSD1306_COMMAND, command})
}

func (oled *SSD1306) writeCommands(commands ...byte) error {
    var command byte
    var err error
    for _, command = range commands {
        _, err = oled.writeCommand(command)
        if err != nil {
            return err
        }
    }
    return nil
}

func (oled *SSD1306) writeData(data byte) (int, error) {
    return oled.i2cConnect.WriteBytes([]byte{SSD1306_DATA, data})
}

func (oled *SSD1306) writeDataBulk(dataBulk []byte) error {
    var data byte
    var err error
    for _, data = range dataBulk {
        _, err = oled.writeData(data)
        if err != nil {
            return err
        }
    }
    return nil
}

// Init i2c Oled Display
func (oled *SSD1306) initDisplay() error {
    err := oled.writeCommands(
        SSD1306_DISPLAY_OFF,
        SSD1306_SET_DISPLAY_CLOCK_DIV, 0x80,
        SSD1306_SET_MULTIPLEX, 0x3F,
        SSD1306_SET_DISPLAY_OFFSET, 0x0,
        SSD1306_SET_START_LINE|0x0,
        SSD1306_CHARGE_PUMP, 0x14,
        SSD1306_MEMORY_MODE, 0x00,
        SSD1306_SEG_REMAP|0x1,
        SSD1306_COM_SCAN_DEC,
        SSD1306_SET_COM_PINS, 0x12,
        SSD1306_SET_CONTRAST, 0xCF,
        SSD1306_SET_PRECHARGE, 0xF1,
        SSD1306_SET_VCOM_DETECT, 0x40,
        SSD1306_DISPLAY_ALL_ON_RESUME,
        SSD1306_NORMAL_DISPLAY,
        SSD1306_DISPLAY_ON,
    )
    if err != nil {
        return err
    }
    err = oled.Clear()
    if err != nil {
        return err
    }

    return nil
}

func (oled *SSD1306) getPageAddress(x, y uint8) PageAddressType {
    return PageAddressType{
        LowerStartColumn: x | 0x1,
        UpperStartColumn: (x | 0xF0) >> 4,
        PageStart:        oled.height / 8,
    }
}

func (oled *SSD1306) Clear() error {
    var err error
    err = oled.writeCommands(
        SSD1306_MEMORY_MODE,
        0x00,
        SSD1306_PAGE_START_ADDRESS_0,
        0x00,
        0x10,
    )
    if err != nil {
        return err
    }

    count := int(oled.width) * int(oled.height) / 8
    for i := 0; i < count; i++ {
        _, err = oled.writeData(0x0)
        if err != nil {
            return err
        }
    }

    return nil
}

func (oled *SSD1306) eachPixel(pixelFunc func(x, y uint8) PixelType) error {
    var x, y uint8
    for x = 0; x < oled.width; x++ {
        for y = 0; y < oled.height; y += 8 {
        }
    }
    return nil
}

func main() {
    oled, err := NewSSD1306(0x3C, 1, 128, 64)
    if err != nil {
        log.Fatal(err)
    }

    oled.writeCommands(0x20, 2)
    oled.writeCommand(2)
    oled.writeCommand(0xB2)
    oled.writeCommand(0x00)
    oled.writeCommand(0x12)

    for _, letter := range font.OledASCIITable {
        oled.writeDataBulk(letter);
        oled.writeDataBulk([]byte{0x0});
        time.Sleep(1 * time.Second)
    }

    time.Sleep(10 * time.Second)

    oled.writeCommand(SSD1306_DISPLAY_OFF)
}
