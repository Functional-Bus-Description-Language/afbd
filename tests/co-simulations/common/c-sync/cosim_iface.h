#ifndef _COSIM_IFACE_H_
#define _COSIM_IFACE_H_

#include <stdint.h>

#include "afbd.h"

typedef uint32_t (*delay_function_t)(void);

void cosim_iface_init(char *wr_fifo_path, char *rd_fifo_path, delay_function_t delay_func);

afbd_iface_t cosim_iface_iface(void);

void cosim_iface_end(int status);

#endif // _COSIM_IFACE_H_
