#include <assert.h>
#include <stdio.h>
#include <time.h>
#include <stdlib.h>

#include "cosim_iface.h"

#include "afbd.h"
#include "main.h"
#define AFBD_IFACE &iface


int main(int argc, char *argv[]) {
	assert(argc == 3);

	afbd_iface_t iface = cosim_iface_iface();

	cosim_iface_init(argv[1], argv[2], NULL);

	srand(time(NULL));
	const uint8_t val = rand() & 0x7F; // 7 bit value

	printf("generated random value: %d\n", val);

	uint8_t cfg;
	uint8_t st;

	afbd_write(main_cfg, val);

	afbd_read(main_cfg, &cfg);
	if (cfg != val) {
		fprintf(stderr, "read wrong value from cfg %d, expecting %d\n", cfg, val);
		cosim_iface_end(1);
	}

	afbd_read(main_st, &st);
	if (st != val) {
		fprintf(stderr, "read wrong value from st %d, expecting %d\n", st, val);
		cosim_iface_end(1);
	}

	cosim_iface_end(0);

	return 0;
}
