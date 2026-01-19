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
	const uint16_t a = rand() & 0xFFFF; // 16 bits value
	const uint16_t b = rand() & 0xFFFF; // 16 bits value

	printf("generated random values: a = %d, b = %d\n", a, b);

	printf("calling add function\n");

	uint32_t err;

	err = afbd_main_add(&iface, a, b);
	if (err) {
		fprintf(stderr, "error calling add function: %d", err);
		cosim_iface_end(1);
	}

	uint32_t result;
	err = afbd_read(main_result, &result);
	if (err) {
		fprintf(stderr, "error reading result: %d", err);
		cosim_iface_end(1);
	}

	if (a + b != result) {
		fprintf(stderr, "wrong result %d, expecting %d\n", result, a + b);
		cosim_iface_end(1);
	}

	cosim_iface_end(0);

	return 0;
}
