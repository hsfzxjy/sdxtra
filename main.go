package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/hsfzxjy/sdxtra/internal/db"
	"github.com/hsfzxjy/sdxtra/internal/ffi"
)

const DSN = "sqlite3://db.sqlite3?cache=shared&_fk=1"

func main() {
	if err := db.Migrate(DSN); err != nil {
		panic(err)
	}
	os.Setenv("CUDA_VISIBLE_DEVICES", "1")
	var models string
	models = "/models"
	go func() {
		for entry := range ffi.GlobalLog() {
			fmt.Printf("[%d] %s", entry.Level, entry.Message)
		}
	}()
	ctx := ffi.CaptureLog0(ffi.GlobalLog(), nil, ffi.SDCtxParams{
		DiffusionModelPath: models + "/flux/flux1-schnell-q4_k.gguf",
		ClipLPath:          models + "/flux/clip_l.safetensors",
		VaePath:            models + "/flux/ae.safetensors",
		T5xxlPath:          models + "/flux/t5xxl_fp16.safetensors",
		NThreads:           -1,
	}.New)
	defer ctx.Free()
	params := ffi.Txt2Img{
		CfgScale:     1.0,
		SampleSteps:  4,
		SampleMethod: ffi.SM_EULER,
		Seed:         40,
		Height:       512,
		Width:        512,
		BatchCount:   1,
		Prompt:       "A photo-realistic illustration of a cute cartoon dog programming and holding a sign saying \"lambdex\".",
	}
	r, err := params.RunOn(ctx)
	if err != nil {
		panic(err)
	}
	defer r.Free()
	f, err := os.Create("output/output.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, r.Image())
	if err != nil {
		panic(err)
	}
	println("done")
}
