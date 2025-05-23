#include <assert.h>
#include <stdio.h>

#include "cosim_iface.h"

#include "afbd.h"
#include "Main.h"
#define AFBD_IFACE &iface


int main(int argc, char *argv[]) {
	assert(argc == 3);

	afbd_iface_t iface = cosim_iface_iface();

	cosim_iface_init(argv[1], argv[2], NULL);

	uint32_t id;
	afbd_read(Main_ID, &id);
	if (id != afbd_Main_ID) {
		fprintf(stderr, "read wrong ID %x, expecting %x\n", id, afbd_Main_ID);
		cosim_iface_end(1);
	}

	uint32_t timestamp;
	afbd_read(Main_TIMESTAMP, &timestamp);
	if (timestamp != afbd_Main_TIMESTAMP) {
		fprintf(stderr, "read wrong TIMESTAMP %x, expecting %x\n", id, afbd_Main_TIMESTAMP);
		cosim_iface_end(1);
	}

	cosim_iface_end(0);

	return 0;
}
