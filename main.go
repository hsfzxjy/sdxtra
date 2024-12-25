package main

// #cgo LDFLAGS: -L../stable-diffusion.cpp/build/bin -lstable-diffusion
// #include "../stable-diffusion.cpp/stable-diffusion.h"
// #include <stdlib.h>
// #include <stdio.h>
// #include <stdbool.h>
// void sd_log_callback(enum sd_log_level_t level, const char* text, void* data) {
// 		printf("[%d] %s", level, text);
// }
import "C"
import (
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"unsafe"
)

type RngType uint32

const (
	STD_DEFAULT_RNG RngType = C.STD_DEFAULT_RNG
	CUDA_RNG        RngType = C.CUDA_RNG
)

type SampleMethod uint32

const (
	EULER_A   SampleMethod = C.EULER_A
	EULER     SampleMethod = C.EULER
	HEUN      SampleMethod = C.HEUN
	DPM2      SampleMethod = C.DPM2
	DPMPP2S_A SampleMethod = C.DPMPP2S_A
	DPMPP2M   SampleMethod = C.DPMPP2M
	DPMPP2Mv2 SampleMethod = C.DPMPP2Mv2
	IPNDM     SampleMethod = C.IPNDM
	IPNDM_V   SampleMethod = C.IPNDM_V
	LCM       SampleMethod = C.LCM
)

type Schedule uint32

const (
	SCHED_DEFAULT  Schedule = C.DEFAULT
	SCHED_DISCRETE Schedule = C.DISCRETE
	KARRAS         Schedule = C.KARRAS
	EXPONENTIAL    Schedule = C.EXPONENTIAL
	AYS            Schedule = C.AYS
	GITS           Schedule = C.GITS
)

type CArena struct {
	ptrs   []unsafe.Pointer
	arenas []*CArena
}

func (a *CArena) CString(s string) *C.char {
	cs := C.CString(s)
	a.ptrs = append(a.ptrs, unsafe.Pointer(cs))
	return cs
}

func (a *CArena) NewArena() *CArena {
	arena := new(CArena)
	a.arenas = append(a.arenas, arena)
	return arena
}

func (a *CArena) Free() {
	ptrs := a.ptrs
	a.ptrs = nil
	for _, ptr := range ptrs {
		C.free(ptr)
	}
	arenas := a.arenas
	a.arenas = nil
	for _, arena := range arenas {
		arena.Free()
	}
}

type SDParams struct {
	ModelPath             string
	ClipLPath             string
	ClipGPath             string
	T5xxlPath             string
	DiffusionModelPath    string
	VaePath               string
	TaesdPath             string
	ControlNetPath        string
	LoraModelDir          string
	EmbedDir              string
	StackedIdEmbedDir     string
	VaeDecodeOnly         bool
	VaeTiling             bool
	FreeParamsImmediately bool
	NThreads              int
	wtype                 int // always = SD_TYPE_COUNT
	RngType               RngType
	Schedule              Schedule
	KeepClipOnCpu         bool
	KeepControlNetCpu     bool
	KeepVaeOnCpu          bool
	DiffusionFlashAttn    bool
}

type SDCtx struct {
	ctx   *C.sd_ctx_t
	arena CArena
}

func NewSDCtx(params SDParams) *SDCtx {
	ctx := new(SDCtx)
	ctx.ctx = C.new_sd_ctx(
		ctx.arena.CString(params.ModelPath),
		ctx.arena.CString(params.ClipLPath),
		ctx.arena.CString(params.ClipGPath),
		ctx.arena.CString(params.T5xxlPath),
		ctx.arena.CString(params.DiffusionModelPath),
		ctx.arena.CString(params.VaePath),
		ctx.arena.CString(params.TaesdPath),
		ctx.arena.CString(params.ControlNetPath),
		ctx.arena.CString(params.LoraModelDir),
		ctx.arena.CString(params.EmbedDir),
		ctx.arena.CString(params.StackedIdEmbedDir),
		C.bool(params.VaeDecodeOnly),
		C.bool(params.VaeTiling),
		C.bool(params.FreeParamsImmediately),
		C.int(params.NThreads),
		C.SD_TYPE_COUNT,
		uint32(params.RngType),
		uint32(params.Schedule),
		C.bool(params.KeepClipOnCpu),
		C.bool(params.KeepControlNetCpu),
		C.bool(params.KeepVaeOnCpu),
		C.bool(params.DiffusionFlashAttn),
	)
	return ctx
}

func (ctx *SDCtx) Free() {
	C.free_sd_ctx(ctx.ctx)
	ctx.arena.Free()
}

type SDImage struct {
	H    int
	W    int
	C    int
	data unsafe.Pointer
}

func _() {
	var _ image.Image = &SDImage{}
}

func newSDImageFromC(cimg *C.sd_image_t) *SDImage {
	sdi := &SDImage{
		H:    int(cimg.height),
		W:    int(cimg.width),
		C:    int(cimg.channel),
		data: unsafe.Pointer(cimg.data),
	}
	C.free(unsafe.Pointer(cimg))
	return sdi
}

func (sdi *SDImage) Free() {
	C.free(unsafe.Pointer(sdi.data))
}

func (sdi *SDImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (sdi *SDImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, sdi.W, sdi.H)
}

func (sdi *SDImage) At(x, y int) color.Color {
	if !(0 <= x && x < sdi.W && 0 <= y && y < sdi.H) {
		return color.RGBA{}
	}
	const C = 3
	i := y*sdi.W*C + x*C
	s := unsafe.Slice((*uint8)(unsafe.Add(sdi.data, i)), C)
	return color.RGBA{s[0], s[1], s[2], 255}
}

type Txt2ImgParams struct {
	Prompt            string
	NegativePrompt    string
	ClipSkip          int
	CfgScale          float32
	Guidance          float32
	Width             int
	Height            int
	SampleMethod      SampleMethod
	SampleSteps       int
	Seed              int64
	BatchCount        int
	ControlCond       *SDImage
	ControlStrength   float32
	StyleStrength     float32
	NormalizeInput    bool
	InputIdImagesPath string
	SkipLayers        []int
	SlgScale          float32
	SkipLayerStart    float32
	SkipLayerEnd      float32
}

func (p Txt2ImgParams) RunOn(ctx *SDCtx) *SDImage {
	arena := ctx.arena.NewArena()
	defer arena.Free()
	img := C.txt2img(
		ctx.ctx,
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
	runtime.KeepAlive(&p)
	return newSDImageFromC(img)
}

func main() {
	os.Setenv("CUDA_VISIBLE_DEVICES", "4")
	C.sd_set_log_callback((*[0]byte)(C.sd_log_callback), nil)
	info := C.sd_get_system_info()
	goinfo := C.GoString(info)
	println(goinfo)
	var models string
	models = "../models"
	models = "/models"
	ctx := NewSDCtx(SDParams{
		DiffusionModelPath: models + "/flux/flux1-schnell-q4_k.gguf",
		ClipLPath:          models + "/flux/clip_l.safetensors",
		VaePath:            models + "/flux/ae.safetensors",
		T5xxlPath:          models + "/flux/t5xxl_fp16.safetensors",
		NThreads:           -1,
	})
	defer ctx.Free()
	params := Txt2ImgParams{
		CfgScale:     1.0,
		SampleSteps:  4,
		SampleMethod: EULER,
		Seed:         100,
		Height:       512,
		Width:        512,
		BatchCount:   1,
		Prompt:       "A photo-realistic illustration of a cute cartoon dog programming and holding a sign saying \"lambdex\".",
	}
	img:=params.RunOn(ctx)
	defer img.Free()
	f, err := os.Create("output/output.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	println("done")
}
