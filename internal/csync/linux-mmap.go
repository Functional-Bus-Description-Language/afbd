package csync

import (
	_ "embed"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
)

//go:embed templates/linux-mmap-iface.h
var mmapHeaderTmplStr string
var mmapHeaderTmpl = template.Must(template.New("C-Sync linux-mmap-iface.h").Parse(mmapHeaderTmplStr))

//go:embed templates/linux-mmap-iface.c
var mmapSourceTmplStr string
var mmapSourceTmpl = template.Must(template.New("C-Sync linux-mmap-iface.c").Parse(mmapSourceTmplStr))

type LinuxMmapFormatters struct {
	BusWidth int64
}

func GenLinuxMmapIface(bus *fn.Block) {
	fmts := LinuxMmapFormatters{
		BusWidth: bus.Width,
	}

	// Generate header file
	hFile, err := os.Create(path.Join(args.CSync.Path, "linux-mmap-iface.h"))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	err = mmapHeaderTmpl.Execute(hFile, fmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	// Generate source file
	cFile, err := os.Create(path.Join(args.CSync.Path, "linux-mmap-iface.c"))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
	err = mmapSourceTmpl.Execute(cFile, fmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
}
