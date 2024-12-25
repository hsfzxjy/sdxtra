package ffi

// #cgo LDFLAGS: -L../../sdcpp/build/shared/bin -lstable-diffusion
// #include "../../sdcpp/stable-diffusion.h"
// #include <stdlib.h>
// #include <stdio.h>
// #include <stdbool.h>
// typedef enum rng_type_t rng_type_t;
// typedef enum sample_method_t sample_method_t;
// typedef enum schedule_t schedule_t;
// typedef enum sd_type_t sd_type_t;
import "C"

type RngTypeT C.rng_type_t

const (
	RNG_STD_DEFAULT_RNG RngTypeT = C.STD_DEFAULT_RNG
	RNG_CUDA_RNG        RngTypeT = C.CUDA_RNG
)

type SampleMethodT C.sample_method_t

const (
	SM_EULER_A   SampleMethodT = C.EULER_A
	SM_EULER     SampleMethodT = C.EULER
	SM_HEUN      SampleMethodT = C.HEUN
	SM_DPM2      SampleMethodT = C.DPM2
	SM_DPMPP2S_A SampleMethodT = C.DPMPP2S_A
	SM_DPMPP2M   SampleMethodT = C.DPMPP2M
	SM_DPMPP2Mv2 SampleMethodT = C.DPMPP2Mv2
	SM_IPNDM     SampleMethodT = C.IPNDM
	SM_IPNDM_V   SampleMethodT = C.IPNDM_V
	SM_LCM       SampleMethodT = C.LCM
)

type ScheduleT C.schedule_t

const (
	SCHED_DEFAULT     ScheduleT = C.DEFAULT
	SCHED_DISCRETE    ScheduleT = C.DISCRETE
	SCHED_KARRAS      ScheduleT = C.KARRAS
	SCHED_EXPONENTIAL ScheduleT = C.EXPONENTIAL
	SCHED_AYS         ScheduleT = C.AYS
	SCHED_GITS        ScheduleT = C.GITS
)

type SDTypeT C.sd_type_t

const (
	SD_TYPE_F32      SDTypeT = C.SD_TYPE_F32
	SD_TYPE_F16      SDTypeT = C.SD_TYPE_F16
	SD_TYPE_Q4_0     SDTypeT = C.SD_TYPE_Q4_0
	SD_TYPE_Q4_1     SDTypeT = C.SD_TYPE_Q4_1
	SD_TYPE_Q5_0     SDTypeT = C.SD_TYPE_Q5_0
	SD_TYPE_Q5_1     SDTypeT = C.SD_TYPE_Q5_1
	SD_TYPE_Q8_0     SDTypeT = C.SD_TYPE_Q8_0
	SD_TYPE_Q8_1     SDTypeT = C.SD_TYPE_Q8_1
	SD_TYPE_Q2_K     SDTypeT = C.SD_TYPE_Q2_K
	SD_TYPE_Q3_K     SDTypeT = C.SD_TYPE_Q3_K
	SD_TYPE_Q4_K     SDTypeT = C.SD_TYPE_Q4_K
	SD_TYPE_Q5_K     SDTypeT = C.SD_TYPE_Q5_K
	SD_TYPE_Q6_K     SDTypeT = C.SD_TYPE_Q6_K
	SD_TYPE_Q8_K     SDTypeT = C.SD_TYPE_Q8_K
	SD_TYPE_IQ2_XXS  SDTypeT = C.SD_TYPE_IQ2_XXS
	SD_TYPE_IQ2_XS   SDTypeT = C.SD_TYPE_IQ2_XS
	SD_TYPE_IQ3_XXS  SDTypeT = C.SD_TYPE_IQ3_XXS
	SD_TYPE_IQ1_S    SDTypeT = C.SD_TYPE_IQ1_S
	SD_TYPE_IQ4_NL   SDTypeT = C.SD_TYPE_IQ4_NL
	SD_TYPE_IQ3_S    SDTypeT = C.SD_TYPE_IQ3_S
	SD_TYPE_IQ2_S    SDTypeT = C.SD_TYPE_IQ2_S
	SD_TYPE_IQ4_XS   SDTypeT = C.SD_TYPE_IQ4_XS
	SD_TYPE_I8       SDTypeT = C.SD_TYPE_I8
	SD_TYPE_I16      SDTypeT = C.SD_TYPE_I16
	SD_TYPE_I32      SDTypeT = C.SD_TYPE_I32
	SD_TYPE_I64      SDTypeT = C.SD_TYPE_I64
	SD_TYPE_F64      SDTypeT = C.SD_TYPE_F64
	SD_TYPE_IQ1_M    SDTypeT = C.SD_TYPE_IQ1_M
	SD_TYPE_BF16     SDTypeT = C.SD_TYPE_BF16
	SD_TYPE_Q4_0_4_4 SDTypeT = C.SD_TYPE_Q4_0_4_4
	SD_TYPE_Q4_0_4_8 SDTypeT = C.SD_TYPE_Q4_0_4_8
	SD_TYPE_Q4_0_8_8 SDTypeT = C.SD_TYPE_Q4_0_8_8
	SD_TYPE_TQ1_0    SDTypeT = C.SD_TYPE_TQ1_0
	SD_TYPE_TQ2_0    SDTypeT = C.SD_TYPE_TQ2_0
	SD_TYPE_COUNT    SDTypeT = C.SD_TYPE_COUNT
)

type sd_image_t = C.sd_image_t
type sd_ctx_t = C.sd_ctx_t

type freeable interface {
	Free()
}
