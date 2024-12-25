package ffi

// #include <stdlib.h>
// #include "../../sdcpp/stable-diffusion.h"
import "C"

type SDCtxParams struct {
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
	RngType               RngTypeT
	Schedule              ScheduleT
	KeepClipOnCpu         bool
	KeepControlNetCpu     bool
	KeepVaeOnCpu          bool
	DiffusionFlashAttn    bool
}

func SDCtxParamsDefault() SDCtxParams {
	return SDCtxParams{
		ModelPath:             "",
		ClipLPath:             "",
		ClipGPath:             "",
		T5xxlPath:             "",
		DiffusionModelPath:    "",
		VaePath:               "",
		TaesdPath:             "",
		ControlNetPath:        "",
		LoraModelDir:          "",
		EmbedDir:              "",
		StackedIdEmbedDir:     "",
		VaeDecodeOnly:         false,
		VaeTiling:             false,
		FreeParamsImmediately: true,
		NThreads:              0,
		wtype:                 int(SD_TYPE_COUNT),
		RngType:               RngTypeT(RNG_STD_DEFAULT_RNG),
		Schedule:              ScheduleT(SM_EULER),
		KeepClipOnCpu:         false,
		KeepControlNetCpu:     false,
		KeepVaeOnCpu:          false,
		DiffusionFlashAttn:    false,
	}
}

func (p SDCtxParams) New() *SDCtx {
	ctx := new(SDCtx)
	arena := &ctx.arena
	cptr := C.new_sd_ctx(
		arena.CString(p.ModelPath),
		arena.CString(p.ClipLPath),
		arena.CString(p.ClipGPath),
		arena.CString(p.T5xxlPath),
		arena.CString(p.DiffusionModelPath),
		arena.CString(p.VaePath),
		arena.CString(p.TaesdPath),
		arena.CString(p.ControlNetPath),
		arena.CString(p.LoraModelDir),
		arena.CString(p.EmbedDir),
		arena.CString(p.StackedIdEmbedDir),
		C.bool(p.VaeDecodeOnly),
		C.bool(p.VaeTiling),
		C.bool(p.FreeParamsImmediately),
		C.int(p.NThreads),
		C.SD_TYPE_COUNT,
		uint32(p.RngType),
		uint32(p.Schedule),
		C.bool(p.KeepClipOnCpu),
		C.bool(p.KeepControlNetCpu),
		C.bool(p.KeepVaeOnCpu),
		C.bool(p.DiffusionFlashAttn),
	)
	if cptr == nil {
		return nil
	}
	ctx.cptr = cptr
	return ctx
}

type SDCtx struct {
	cptr  *sd_ctx_t
	arena cgoArena
}

func (ctx *SDCtx) Free() {
	C.free_sd_ctx(ctx.cptr)
	ctx.arena.Free()
}
