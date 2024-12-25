#ifndef SDXTRA_INTERNAL_FFI_LOG_H
#define SDXTRA_INTERNAL_FFI_LOG_H
#include "../../sdcpp/stable-diffusion.h"
#include <pthread.h>
typedef enum sd_log_level_t sd_log_level_t;
extern void goHandleLog(sd_log_level_t level, char* text, uintptr_t data, pthread_t threadId);
extern void handleLog(sd_log_level_t level, char* text, uintptr_t data);
#endif