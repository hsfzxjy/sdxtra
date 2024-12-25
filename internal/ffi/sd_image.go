package ffi

// #include <stdlib.h>
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

type SDRGBImage struct {
	Width  int
	Height int
	data   unsafe.Pointer
}

func newSDRGBImage(cimg *sd_image_t) (*SDRGBImage, error) {
	defer C.free(unsafe.Pointer(cimg))
	if cimg.channel != 3 {
		return nil, fmt.Errorf("unsupported channel count: %d", cimg.channel)
	}
	sdi := &SDRGBImage{
		Width:  int(cimg.width),
		Height: int(cimg.height),
		data:   unsafe.Pointer(cimg.data),
	}
	return sdi, nil
}

func (sdi *SDRGBImage) Free() {
	C.free(unsafe.Pointer(sdi.data))
}

func (sdi *SDRGBImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (sdi *SDRGBImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, sdi.Width, sdi.Height)
}

func (sdi *SDRGBImage) At(x, y int) color.Color {
	if !(0 <= x && x < sdi.Width && 0 <= y && y < sdi.Height) {
		return color.RGBA{}
	}
	const C = 3
	i := y*sdi.Width*C + x*C
	s := unsafe.Slice((*uint8)(unsafe.Add(sdi.data, i)), C)
	return color.RGBA{s[0], s[1], s[2], 255}
}
