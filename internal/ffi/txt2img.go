package ffi

// #include <stdlib.h>
// #include "../../sdcpp/stable-diffusion.h"
import "C"

import (
	"fmt"
	"image"
	"unsafe"
)

type Txt2Img struct {
	Prompt            string
	NegativePrompt    string
	ClipSkip          int
	CfgScale          float32
	Guidance          float32
	Width             int
	Height            int
	SampleMethod      SampleMethodT
	SampleSteps       int
	Seed              int64
	BatchCount        int
	ControlCond       *SDRGBImage
	ControlStrength   float32
	StyleStrength     float32
	NormalizeInput    bool
	InputIdImagesPath string
	SkipLayers        []int
	SlgScale          float32
	SkipLayerStart    float32
	SkipLayerEnd      float32
}

func Txt2ImgDefault() Txt2Img {
	return Txt2Img{
		ClipSkip:        0,
		CfgScale:        1.0,
		Guidance:        0.0,
		Width:           512,
		Height:          512,
		SampleMethod:    SM_EULER,
		SampleSteps:     0,
		Seed:            0,
		BatchCount:      1,
		ControlStrength: 0.0,
		StyleStrength:   0.0,
		NormalizeInput:  false,
		SkipLayers:      nil,
		SlgScale:        0.0,
		SkipLayerStart:  0.0,
		SkipLayerEnd:    0.0,
	}
}

type Txt2ImgResult struct {
	image *SDRGBImage
}

func (r *Txt2ImgResult) Free() {
	r.image.Free()
}

func (r *Txt2ImgResult) Image() image.Image {
	return r.image
}

func (p Txt2Img) RunOn(ctx *SDCtx) (*Txt2ImgResult, error) {
	arena := ctx.arena.NewArena()
	defer arena.Free()
	cimg := p.runOn(ctx, arena)
	if cimg == nil {
		return nil, fmt.Errorf("txt2img failed")
	}
	img, err := newSDRGBImage(cimg)
	if err != nil {
		return nil, err
	}
	return &Txt2ImgResult{
		image: img,
	}, nil
}

func (p *Txt2Img) runOn(ctx *SDCtx, arena *cgoArena) *sd_image_t {
	return C.txt2img(
		ctx.cptr,
		arena.CString(p.Prompt),
		arena.CString(p.NegativePrompt),
		C.int(p.ClipSkip),
		C.float(p.CfgScale),
		C.float(p.Guidance),
		C.int(p.Width),
		C.int(p.Height),
		uint32(p.SampleMethod),
		C.int(p.SampleSteps),
		C.int64_t(p.Seed),
		C.int(p.BatchCount),
		nil, // TODO
		C.float(p.ControlStrength),
		C.float(p.StyleStrength),
		C.bool(p.NormalizeInput),
		arena.CString(p.InputIdImagesPath),
		(*C.int)(unsafe.Pointer(unsafe.SliceData(p.SkipLayers))),
		C.size_t(len(p.SkipLayers)),
		C.float(p.SlgScale),
		C.float(p.SkipLayerStart),
		C.float(p.SkipLayerEnd),
	)
}
