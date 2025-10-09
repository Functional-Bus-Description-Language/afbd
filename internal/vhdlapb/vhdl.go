package vhdlapb

import (
	"log"
	"os"
	"sync"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
	"github.com/Functional-Bus-Description-Language/afbd/internal/utils"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

var busWidth int64

func Generate(bus *fn.Block, pkgsConsts map[string]*pkg.Package) {
	busWidth = bus.Width

	err := os.MkdirAll(args.VhdlApb.Path, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate vhdl-apb: %v", err)
	}

	blocks := utils.CollectBlocks(bus, nil, []string{})
	utils.ResolveBlockNameConflicts(blocks)

	var wg sync.WaitGroup
	defer wg.Wait()

	genAPBPackage(pkgsConsts)

	for _, b := range blocks {
		wg.Add(1)
		go genBlock(b, &wg)
	}
}
