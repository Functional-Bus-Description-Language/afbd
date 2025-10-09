package csync

import (
	_ "embed"
	"log"
	"os"
	"path"
	"sync"
	"text/template"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
	"github.com/Functional-Bus-Description-Language/afbd/internal/c"
	"github.com/Functional-Bus-Description-Language/afbd/internal/utils"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

var busWidth int64

var addrType c.Type
var readType c.Type
var writeType c.Type

//go:embed templates/afbd.h
var afbdHeaderTmplStr string
var afbdHeaderTmpl = template.Must(template.New("C-Sync afbd.h").Parse(afbdHeaderTmplStr))

type afbdHeaderFormatters struct {
	AddrType  string
	ReadType  string
	WriteType string
}

func Generate(bus *fn.Block, pkgsConsts map[string]*pkg.Package) {
	busWidth = bus.Width

	err := os.MkdirAll(args.CSync.Path, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	hFile, err := os.Create(path.Join(args.CSync.Path, "afbd.h"))
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	addrType = c.SizeToAddrType(bus.Sizes.BlockAligned)
	readType = c.WidthToReadType(bus.Width)
	writeType = c.WidthToWriteType(bus.Width)

	hFmts := afbdHeaderFormatters{
		AddrType:  addrType.String(),
		ReadType:  readType.String(),
		WriteType: writeType.String(),
	}

	err = afbdHeaderTmpl.Execute(hFile, hFmts)
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}

	blocks := utils.CollectBlocks(bus, nil, []string{})
	utils.ResolveBlockNameConflicts(blocks)

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, b := range blocks {
		wg.Add(1)
		go genBlock(b, &wg)
	}

	err = hFile.Close()
	if err != nil {
		log.Fatalf("generate C-Sync: %v", err)
	}
}
