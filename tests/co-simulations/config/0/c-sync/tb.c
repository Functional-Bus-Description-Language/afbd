#include <assert.h>
#include <stdio.h>
#include <time.h>
#include <stdlib.h>

#include "cosim_iface.h"

#include "afbd.h"
#include "Main.h"
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

	afbd_write(Main_Cfg, val);

	afbd_read(Main_Cfg, &cfg);
	if (cfg != val) {
		fprintf(stderr, "read wrong value from Cfg %d, expecting %d\n", cfg, val);
		cosim_iface_end(1);
	}

	afbd_read(Main_Cfg, &st);
	if (st != val) {
		fprintf(stderr, "read wrong value from St %d, expecting %d\n", st, val);
		cosim_iface_end(1);
	}

	cosim_iface_end(0);

	return 0;
}
