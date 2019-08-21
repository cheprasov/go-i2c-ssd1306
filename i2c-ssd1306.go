//package i2c_oled
package main

import (
    "image"
    "image/color"
    "image/png"
    "io/ioutil"
    "log"
    "strings"
    "time"

    "./font"
    "./helpers"
    "github.com/cheprasov/go-i2c"
)

const (
    SSD1306_CODC_COMMAND    = 0x80
    SSD1306_CODC_DATA       = 0xC0
    SSD1306_CODC_MULTI_DATA = 0x40
)

const (
    SSD1306_PAGE_SIZE = 8
)

const (
    SSD1306_CMD_I2C_ADDRESS           = 0x3C // 011110+SA0+RW - 0x3C or 0x3D
    SSD1306_CMD_SET_CONTRAST          = 0x81
    SSD1306_CMD_DISPLAY_ALL_ON_RESUME = 0xA4
    SSD1306_CMD_DISPLAY_ALL_ON        = 0xA5
    SSD1306_CMD_NORMAL_DISPLAY        = 0xA6
    SSD1306_CMD_INVERT_DISPLAY        = 0xA7
    SSD1306_CMD_DISPLAY_OFF           = 0xAE
    SSD1306_CMD_DISPLAY_ON            = 0xAF
    SSD1306_CMD_SET_DISPLAY_OFFSET    = 0xD3
    SSD1306_CMD_SET_COM_PINS          = 0xDA
    SSD1306_CMD_SET_VCOM_DETECT       = 0xDB
    SSD1306_CMD_SET_DISPLAY_CLOCK_DIV = 0xD5
    SSD1306_CMD_SET_PRECHARGE         = 0xD9
    SSD1306_CMD_SET_MULTIPLEX         = 0xA8
    SSD1306_CMD_SET_LOW_COLUMN        = 0x00
    SSD1306_CMD_SET_HIGH_COLUMN       = 0x10
    SSD1306_CMD_SET_START_LINE        = 0x40
    SSD1306_CMD_MEMORY_MODE           = 0x20
    SSD1306_CMD_SET_COLUMN_ADDRESS    = 0x21
    SSD1306_CMD_SET_PAGE_ADDRESS      = 0x22
    SSD1306_CMD_COM_SCAN_INC          = 0xC0
    SSD1306_CMD_COM_SCAN_DEC          = 0xC8
    SSD1306_CMD_SEG_REMAP             = 0xA0
    SSD1306_CMD_CHARGE_PUMP           = 0x8D
    SSD1306_CMD_EXTERNAL_VCC          = 0x1
    SSD1306_CMD_SWITCH_CAP_VCC        = 0x2

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

    SSD1306_MEMORY_MODE_HORIZONTAL_ADDRESSING = 0x0;
    SSD1306_MEMORY_MODE_VERTICAL_ADDRESSING   = 0x1;
    SSD1306_MEMORY_MODE_PAGE_ADDRESSING       = 0x2;
)

type PixelType uint8;

type PageAddressType struct {
    PageStart          uint8
    PageAddressStart   uint8
    PageAddressEnd     uint8
    LowerStartColumn   uint8
    UpperStartColumn   uint8
    ColumnAddressStart uint8
    ColumnAddressEnd   uint8
}

type SSD1306 struct {
    i2cAddress uint8
    i2cBus     uint8
    i2cConnect *i2c.I2C
    width      uint8
    height     uint8
    pagesCount uint8
}

// Create and init new instance of SSD1306 Oled Display
func NewSSD1306(i2cAddress, i2cBus, width, height uint8) (*SSD1306, error) {
    var err error

    oled := &SSD1306{
        i2cAddress: i2cAddress,
        i2cBus:     i2cBus,
        width:      width,
        height:     height,
        pagesCount: height / SSD1306_PAGE_SIZE,
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

func (oled *SSD1306) GetWidth() uint8 {
    return oled.width
}

func (oled *SSD1306) GetHeight() uint8 {
    return oled.height
}

func (oled *SSD1306) GetPagesCount() uint8 {
    return oled.pagesCount
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

// Init i2c connection
func (oled *SSD1306) Close() error {
    return oled.i2cConnect.Close()
}

func (oled *SSD1306) writeCommand(command byte) (int, error) {
    return oled.i2cConnect.WriteBytes([]byte{SSD1306_CODC_COMMAND, command})
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
    return oled.i2cConnect.WriteBytes([]byte{SSD1306_CODC_DATA, data})
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
        SSD1306_CMD_DISPLAY_OFF,
        SSD1306_CMD_SET_DISPLAY_CLOCK_DIV, 0x80,
        SSD1306_CMD_SET_MULTIPLEX, 0x3F,
        SSD1306_CMD_SET_DISPLAY_OFFSET, 0x0,
        SSD1306_CMD_SET_START_LINE|0x0,
        SSD1306_CMD_CHARGE_PUMP, 0x14,
        SSD1306_CMD_MEMORY_MODE, 0x00,
        SSD1306_CMD_SEG_REMAP|0x1,
        SSD1306_CMD_COM_SCAN_DEC,
        SSD1306_CMD_SET_COM_PINS, 0x12,
        SSD1306_CMD_SET_CONTRAST, 0xCF,
        SSD1306_CMD_SET_PRECHARGE, 0xF1,
        SSD1306_CMD_SET_VCOM_DETECT, 0x40,
        SSD1306_CMD_DISPLAY_ALL_ON_RESUME,
        SSD1306_CMD_NORMAL_DISPLAY,
        SSD1306_CMD_DISPLAY_ON,
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

func (oled *SSD1306) getPageAddress(page, offset, pages, width uint8) *PageAddressType {
    return &PageAddressType{
        PageStart:          SSD1306_PAGE_START_ADDRESS_0 + page,
        PageAddressStart:   page,
        PageAddressEnd:     helpers.IfUint8(pages == 0, 0, helpers.MinUint8(page+pages, oled.pagesCount-1)),
        LowerStartColumn:   offset & 0xF,
        UpperStartColumn:   (offset & 0xF0) >> 4,
        ColumnAddressStart: offset,
        ColumnAddressEnd:   helpers.IfUint8(width == 0, 0, helpers.MinUint8(offset+width, oled.width-1)),
    }
}

func (oled *SSD1306) setPageAddress(pageAddress *PageAddressType) error {
    var err error

    switch true {
    case pageAddress.PageAddressStart < pageAddress.PageAddressEnd:
        err = oled.writeCommands(
            SSD1306_CMD_MEMORY_MODE, SSD1306_MEMORY_MODE_HORIZONTAL_ADDRESSING,
            SSD1306_CMD_SET_PAGE_ADDRESS, pageAddress.PageAddressStart, pageAddress.PageAddressEnd,
        )
    default:
        err = oled.writeCommands(
            SSD1306_CMD_MEMORY_MODE, SSD1306_MEMORY_MODE_HORIZONTAL_ADDRESSING,
            pageAddress.PageStart,
        )
    }

    if err != nil {
        return err
    }

    switch true {
    case pageAddress.ColumnAddressStart < pageAddress.ColumnAddressEnd:
        err = oled.writeCommands(
            SSD1306_CMD_SET_COLUMN_ADDRESS, pageAddress.ColumnAddressStart, pageAddress.ColumnAddressEnd,
        )
    default:
        err = oled.writeCommands(pageAddress.LowerStartColumn, pageAddress.LowerStartColumn)
    }

    if err != nil {
        return err
    }

    return nil
}

func (oled *SSD1306) Clear() error {
    var err error

    err = oled.setPageAddress(oled.getPageAddress(0, 0, 0, 0));
    if err != nil {
        return err
    }

    count := int(oled.width) * int(oled.pagesCount)
    for i := 0; i < count; i++ {
        _, err = oled.writeData(0x00)
        if err != nil {
            return err
        }
    }

    return nil
}

func (oled *SSD1306) setPageCursor(page, offset uint8) error {
    return oled.writeCommands(
        SSD1306_CMD_MEMORY_MODE, SSD1306_MEMORY_MODE_HORIZONTAL_ADDRESSING,
        SSD1306_PAGE_START_ADDRESS_0+page,
        SSD1306_CMD_SET_LOW_COLUMN+(offset&0x0F),
        SSD1306_CMD_SET_HIGH_COLUMN+((offset&0xF0)>>4),
    )
}

func (oled *SSD1306) PrintText(text string, row, offset uint8) error {
    var err error
    err = oled.setPageCursor(row, offset)
    if err != nil {
        return err
    }

    if text == "" {
        return nil
    }

    var r rune
    var charBytes []byte
    var isOk bool
    for _, r = range []rune(text) {
        charBytes, isOk = font.Chars[r]
        if isOk == false {
            charBytes = font.CharUnknown
        }
        err = oled.writeDataBulk(charBytes)
        if err != nil {
            return err
        }
    }

    return nil
}

func (oled *SSD1306) DrawImage(imgPointer *image.Image, page, offset uint8) error {
    var err error
    var img = *imgPointer;
    imgWidth := uint8(img.Bounds().Max.X - img.Bounds().Min.X)
    imgHeight := uint8(img.Bounds().Max.Y - img.Bounds().Min.Y)

    if imgWidth > oled.width {
        imgWidth = oled.width
    }
    if imgWidth + offset > oled.width {
        imgWidth = oled.width - offset
    }

    imgPages := uint8(imgHeight / SSD1306_PAGE_SIZE)
    if imgHeight%SSD1306_PAGE_SIZE != 0 {
        imgPages += 1
    }

    err = oled.setPageAddress(oled.getPageAddress(page, offset, imgPages, imgWidth));
    if err != nil {
        return err
    }

    var currentPage, pixels, pageY, x, y uint8;
    for currentPage = 0; currentPage < imgPages; currentPage++ {
        for x = 0; x < imgWidth; x++ {
            pixels = 0
            for pageY = 0; pageY < SSD1306_PAGE_SIZE; pageY++ {
                y = currentPage*SSD1306_PAGE_SIZE + pageY;
                if y >= imgHeight {
                    continue
                }
                clr := color.GrayModel.Convert(img.At(int(x), int(y))).(color.Gray)
                if clr.Y < 127 {
                    continue
                }
                pixels = pixels | (0x01 << pageY)
            }
            _, err := oled.writeData(pixels)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func gopherPNG() ([]byte, error) {
    return ioutil.ReadFile("./imgs/marilyn_monroe.png")
}

func draw(oled *SSD1306) {
    // This example uses png.Decode which can only decode PNG images.
    // Consider using the general image.Decode as it can sniff and decode any registered image format.
    file, err := gopherPNG()
    if err != nil {
        log.Fatal(err)
    }

    img, err := png.Decode(strings.NewReader(string(file)))
    if err != nil {
        log.Fatal(err)
    }

    err = oled.DrawImage(&img, 4, 32)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    oled, err := NewSSD1306(0x3C, 1, 128, 64)
    if err != nil {
        log.Fatal(err)
    }
    defer oled.Close()

    oled.PrintText("Init", 4, 54)

    time.Sleep(1 * time.Second)

    draw(oled)

    time.Sleep(10 * time.Second)

    oled.writeCommand(SSD1306_CMD_DISPLAY_OFF)
}
