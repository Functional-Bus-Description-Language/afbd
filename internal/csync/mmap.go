package csync

import (
	_ "embed"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
	"github.com/Functional-Bus-Description-Language/afbd/internal/c"
)

//go:embed templates/mmap-iface.h
var mmapHeaderTmplStr string
var mmapHeaderTmpl = template.Must(template.New("C-Sync mmap-iface.h").Parse(mmapHeaderTmplStr))

//go:embed templates/mmap-iface.c
var mmapSourceTmplStr string
var mmapSourceTmpl = template.Must(template.New("C-Sync mmap-iface.c").Parse(mmapSourceTmplStr))

type mmapFormatters struct {
	AddrType       string
	ReadType       string
	WriteType      string
	WordByteShift  int64
	LinuxReadFunc  string
	LinuxWriteFunc string
}

func GenMmapIface() {
	fmts := mmapFormatters{
		AddrType:       addrType.String(),
		ReadType:       readType.String(),
		WriteType:      writeType.String(),
		WordByteShift:  c.WidthToWordByteShift(busWidth),
		LinuxReadFunc:  c.WidthToLinuxReadFunc(busWidth),
		LinuxWriteFunc: c.WidthToLinuxWriteFunc(busWidth),
	}

	// Generate header file
	hFile, err := os.Create(path.Join(args.CSync.Path, "mmap-iface.h"))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	err = mmapHeaderTmpl.Execute(hFile, fmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	// Generate source file
	cFile, err := os.Create(path.Join(args.CSync.Path, "mmap-iface.c"))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	err = mmapSourceTmpl.Execute(cFile, fmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
}
