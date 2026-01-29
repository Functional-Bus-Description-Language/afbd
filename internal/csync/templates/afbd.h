#ifndef _AFBD_AFBD_H_
#define _AFBD_AFBD_H_

#ifdef __KERNEL__
	#include <linux/types.h>
#else
	#include <stddef.h>
	#include <stdint.h>
#endif

typedef struct {
	int (*read)(const {{.AddrType}} addr, {{.ReadType}} const data);
	int (*write)(const {{.AddrType}} addr, const {{.WriteType}} data);
	int (*readb)(const {{.AddrType}} addr, {{.ReadType}} buf, size_t count);
	int (*writeb)(const {{.AddrType}} addr, const {{.WriteType}} * buf, size_t count);
} afbd_iface_t;

#define afbd_read(elem, data) (afbd_ ## elem ## _read(AFBD_IFACE, data))
#define afbd_write(elem, data) (afbd_ ## elem ## _write(AFBD_IFACE, data))

#ifdef AFBD_SHORT_MACROS
	#undef afbd_read
	#undef afbd_write
	#define read(elem, data) (afbd_ ## elem ## _read(AFBD_IFACE, data))
	#define write(elem, data) (afbd_ ## elem ## _write(AFBD_IFACE, data))
#endif

#endif // _AFBD_AFBD_H_
