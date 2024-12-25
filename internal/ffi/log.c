#include "log.h"

void handleLog(sd_log_level_t level, char* text, uintptr_t data) {
   goHandleLog(level, text, data, pthread_self());
}
