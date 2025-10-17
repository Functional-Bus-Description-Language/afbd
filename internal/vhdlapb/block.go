package vhdlapb

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"sync"
	"text/template"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
	"github.com/Functional-Bus-Description-Language/afbd/internal/utils"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/cnst"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

//go:embed templates/block.vhd
var blockEntityTmplStr string
var blockEntityTmpl = template.Must(template.New("vhdl-apb entity").Parse(blockEntityTmplStr))

type BlockEntityFormatters struct {
	BusWidth   int64
	EntityName string

	MasterCount          int64
	RegCount             int64
	InternalAddrBitCount int64
	SubblockCount        int64

	// Things going to package.
	Constants   string
	ProcTypes   string
	StreamTypes string

	EntitySubblockPorts   string
	EntityFunctionalPorts string

	CrossbarSubblockPortsIn  string
	CrossbarSubblockPortsOut string

	SignalDeclarations string
	AddressValues      string
	MaskValues         string

	RegistersAccess RegisterMap

	ProcsCallsClear string
	ProcsCallsSet   string
	ProcsExitsClear string
	ProcsExitsSet   string

	StreamsStrobesClear string
	StreamsStrobesSet   string

	DefaultValues string

	CombinationalProcesses string
}

func genBlock(blk utils.Block, wg *sync.WaitGroup) {
	defer wg.Done()

	intAddrBitCount := int64(0)
	if blk.Block.Sizes.Own > 1 {
		intAddrBitCount = int64(math.Ceil(math.Log2(float64(blk.Block.Sizes.Own))))
	}

	fmts := BlockEntityFormatters{
		BusWidth:             busWidth,
		EntityName:           blk.Name,
		MasterCount:          blk.Block.Masters,
		RegCount:             blk.Block.Sizes.Own,
		InternalAddrBitCount: intAddrBitCount,
		AddressValues:        fmt.Sprintf("0 => \"%032b\"", 0),
		RegistersAccess:      make(RegisterMap),
	}

	addrBitsCount := int(math.Log2(float64(blk.Block.Sizes.BlockAligned)))

	mask := 0
	if len(blk.Block.Subblocks) > 0 {
		mask = ((1 << addrBitsCount) - 1) ^ ((1 << fmts.InternalAddrBitCount) - 1)
	}
	fmts.MaskValues = fmt.Sprintf("0 => \"%032b\"", mask<<2)

	genConsts(&blk.Block.Consts, &fmts)

	for _, sb := range blk.Block.Subblocks {
		genSubblock(sb, blk.Block.StartAddr(), addrBitsCount, &fmts)
	}

	for _, proc := range blk.Block.Procs {
		genProc(proc, &fmts)
	}

	for _, stream := range blk.Block.Streams {
		genStream(stream, &fmts)
	}

	for _, st := range blk.Block.Statics {
		genStatic(st, &fmts)
	}

	for _, st := range blk.Block.Statuses {
		genStatus(st, &fmts)
	}

	for _, cfg := range blk.Block.Configs {
		genConfig(cfg, &fmts)
	}

	for _, mask := range blk.Block.Masks {
		genMask(mask, &fmts)
	}

	filePath := path.Join(args.VhdlApb.Path, (blk.Name + ".vhd"))
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("generate vhdl-apb: %v", err)
	}

	err = blockEntityTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate vhdl-apb: %v", err)
	}

	addGeneratedFile(filePath)

	err = f.Close()
	if err != nil {
		log.Fatalf("generate vhdl-apb: %v", err)
	}
}

func genSubblock(
	sb *fn.Block,
	superBlockAddrStart int64,
	superBlockAddrBitsCount int,
	fmts *BlockEntityFormatters,
) {
	initSubblockCount := fmts.SubblockCount

	fmts.EntitySubblockPorts += fmt.Sprintf(`;
  %s_apb_reqs_o : out apb.requester_out_array_t(0 to %d);
  %[1]s_apb_reqs_i : in  apb.requester_in_array_t(0 to %[2]d)`,
		sb.Name, sb.Count-1,
	)

	if sb.Count == 1 {
		fmts.CrossbarSubblockPortsIn += fmt.Sprintf(
			"\n  reqs_i(%d) => %s_apb_reqs_i(0),",
			initSubblockCount+1, sb.Name,
		)

		fmts.CrossbarSubblockPortsOut += fmt.Sprintf(
			",\n  reqs_o(%d) => %s_apb_reqs_o(0)",
			initSubblockCount+1, sb.Name,
		)
	} else {
		lowerBound := initSubblockCount + 1
		upperBound := lowerBound + sb.Count - 1

		s := fmt.Sprintf("\n  reqs_i(%d to %d) => %s_apb_reqs_i,", lowerBound, upperBound, sb.Name)
		fmts.CrossbarSubblockPortsIn += s

		s = fmt.Sprintf(",\n  reqs_o(%d to %d) => %s_apb_reqs_o", lowerBound, upperBound, sb.Name)
		fmts.CrossbarSubblockPortsOut += s
	}

	subblockAddr := sb.StartAddr() - superBlockAddrStart
	for range sb.Count {
		fmts.SubblockCount += 1

		fmts.AddressValues += fmt.Sprintf(
			", %d => \"%032b\"", fmts.SubblockCount, subblockAddr<<2,
		)

		mask := ((1 << superBlockAddrBitsCount) - 1) ^ ((1 << int(math.Log2(float64(sb.Sizes.BlockAligned)))) - 1)
		fmts.MaskValues += fmt.Sprintf(
			", %d => \"%032b\"", fmts.SubblockCount, mask<<2,
		)

		subblockAddr += sb.Sizes.BlockAligned
	}
}

func genConsts(cc *cnst.Container, fmts *BlockEntityFormatters) {
	s := ""

	for name, b := range cc.Bools {
		s += fmt.Sprintf("constant %s : boolean := %t;\n", name, b)
	}
	for name, list := range cc.BoolLists {
		s += fmt.Sprintf("constant %s : boolean_vector(0 to %d) := (", name, len(list)-1)
		for i, b := range list {
			s += fmt.Sprintf("%d => %t, ", i, b)
		}
		s = s[:len(s)-2]
		s += ");\n"
	}
	for name, i := range cc.Ints {
		s += fmt.Sprintf("constant %s : int64 := signed'(x\"%016x\");\n", name, i)
	}
	for name, list := range cc.IntLists {
		s += fmt.Sprintf("constant %s : int64_vector(0 to %d) := (", name, len(list)-1)
		for i, v := range list {
			s += fmt.Sprintf("%d => signed'(x\"%016x\"), ", i, v)
		}
		s = s[:len(s)-2]
		s += ");\n"
	}
	for name, str := range cc.Strings {
		s += fmt.Sprintf("constant %s : string := %q;\n", name, str)
	}

	fmts.Constants += s
}
