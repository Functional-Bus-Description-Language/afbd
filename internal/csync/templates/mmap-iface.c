#if defined(__KERNEL__) && defined(__linux__)
	#include <linux/io.h>
#endif

#include "mmap-iface.h"

static int mmap_iface_read(afbd_iface_t * const iface, const {{.AddrType}} addr, {{.ReadType}} const data)
{
	const size_t offset = addr << {{.WordByteShift}};

#if defined(__KERNEL__) && defined(__linux__)
	*data = {{.LinuxReadFunc}}(iface->data + offset);
	return 0;
#else
	#error "unimplemented"
#endif
}

int mmap_iface_write(afbd_iface_t * const iface, const {{.AddrType}} addr, const {{.WriteType}} data)
{
	const size_t offset = addr << {{.WordByteShift}};

#if defined(__KERNEL__) && defined(__linux__)
	{{.LinuxWriteFunc}}(data, iface->data + offset);
	return 0;
#else
	#error "unimplemented"
#endif
}

int mmap_iface_readb(afbd_iface_t * const iface, const {{.AddrType}} addr, {{.ReadType}} buf, size_t count)
{
	for (size_t i = 0; i < count; i++)
		iface->read(iface, addr + i, &buf[i]);

	return count;
}

int mmap_iface_writeb(afbd_iface_t * const iface, const {{.AddrType}} addr, const {{.WriteType}} * buf, size_t count)
{
	for (size_t i = 0; i < count; i++)
		iface->write(iface, addr + i, buf[i]);

	return count;
}

afbd_iface_t afbd_mmap_iface(void *mem)
{
	return (afbd_iface_t){
		read: mmap_iface_read,
		write: mmap_iface_write,
		readb: mmap_iface_readb,
		writeb: mmap_iface_writeb,
		data: mem
	};
}
