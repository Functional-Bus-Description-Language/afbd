package csync

import (
	"fmt"
	"strings"

	"github.com/Functional-Bus-Description-Language/afbd/internal/c"
	_ "github.com/Functional-Bus-Description-Language/afbd/internal/utils"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
)

func genProc(p *fn.Proc, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	sig := genProcSignature(p, blk, hFmts)

	hFmts.Code += "\n" + sig + ";\n"

	cFmts.Code += fmt.Sprintf("\n%s\n{\n", sig)
	if len(p.Params) == 0 && len(p.Returns) == 0 {
		callAddr := *p.CallAddr
		if !args.CSync.OffsetAddr {
			callAddr += blk.StartAddr()
		}
		cFmts.Code += fmt.Sprintf("\treturn iface->write(iface, %d, 0);\n", callAddr)
	}

	if len(p.Params) > 0 {
		genProcParamsAccess(p, blk, cFmts)
	}

	if len(p.Returns) > 0 {
		genProcReturnsAccess(p, blk, cFmts)
	}

	cFmts.Code += "};\n"
}

func genProcSignature(p *fn.Proc, blk *fn.Block, hFmts *BlockHFormatters) string {
	prefix := "int afbd_" + hFmts.BlockName + "_" + p.Name

	params := strings.Builder{}
	params.WriteString("afbd_iface_t * const iface")

	for _, p := range p.Params {
		params.WriteString(
			", const " + c.WidthToWriteType(p.Width).String() + " " + p.Name,
		)
	}

	for _, r := range p.Returns {
		params.WriteString(
			", " + c.WidthToReadType(r.Width).String() + " const " + r.Name,
		)
	}

	return prefix + "(" + params.String() + ")"
}

func genProcParamsAccess(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.ParamsBufSize() == 1 {
		genProcParamsAccessSingleWrite(p, blk, cFmts)
	} else {
		genProcParamsAccessBlockWrite(p, blk, cFmts)
	}
}

func genProcParamsAccessSingleWrite(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.Delay == nil && len(p.Returns) == 0 {
		genProcParamsAccessSingleWriteNoDelayNoReturns(p, blk, cFmts)
	} else {
		panic("unimplemented")
	}
}

func genProcParamsAccessSingleWriteNoDelayNoReturns(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	callAddr := *p.CallAddr
	if !args.CSync.OffsetAddr {
		callAddr += blk.StartAddr()
	}

	cFmts.Code += fmt.Sprintf("\treturn iface->write(iface, %d, ", callAddr)
	for i, p := range p.Params {
		if i != 0 {
			cFmts.Code += " | "
		}

		switch acs := p.Access; acs.Type {
		case "SingleOneReg":
			cFmts.Code += fmt.Sprintf("%s << %d", p.Name, acs.StartBit)
		default:
			panic("unimplemented")
		}
	}
	cFmts.Code += ");\n"
}

func genProcParamsAccessBlockWrite(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.Delay == nil && len(p.Returns) == 0 {
		genProcParamsAccessBlockWriteNoDelayNoReturns(p, blk, cFmts)
	} else {
		panic("unimplemented")
	}
}

func genProcParamsAccessBlockWriteNoDelayNoReturns(proc *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	cFmts.Code += fmt.Sprintf("\t%s buf[%d] = {0};\n\n", c.WidthToWriteType(blk.Width), proc.ParamsBufSize())

	for _, param := range proc.Params {
		switch acs := param.Access; acs.Type {
		case "SingleOneReg":
			cFmts.Code += fmt.Sprintf(
				"\tbuf[%d] |= %s << %d;\n",
				acs.StartAddr-proc.ParamsStartAddr(), param.Name, acs.StartBit,
			)
		default:
			panic("unimplemented")
		}
	}

	startAddr := proc.ParamsStartAddr()
	if !args.CSync.OffsetAddr {
		startAddr += blk.StartAddr()
	}

	cFmts.Code += fmt.Sprintf(
		"\n\treturn iface->writeb(iface, %d, buf, %d);\n",
		startAddr, proc.ParamsBufSize(),
	)
}

func genProcReturnsAccess(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	if p.ReturnsBufSize() == 1 {
		genProcReturnsAccessSingleRead(p, blk, cFmts)
	} else {
		genProcReturnsAccessBlockRead(p, blk, cFmts)
	}
}

func genProcReturnsAccessSingleRead(proc *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	cFmts.Code += fmt.Sprintf("\t%s _rdata;\n", c.WidthToWriteType(blk.Width))

	exitAddr := *proc.ExitAddr
	if !args.CSync.OffsetAddr {
		exitAddr += blk.StartAddr()
	}

	cFmts.Code += fmt.Sprintf("\tconst int err = iface->read(%d, &_rdata);\n", exitAddr)
	cFmts.Code += "\tif (err)\n\t\t return err;\n"

	for _, r := range proc.Returns {
		switch acs := r.Access; acs.Type {
		case "SingleOneReg":
			cFmts.Code += fmt.Sprintf(
				"\t*%s = (_rdata >> %d) & 0x%X;\n",
				r.Name, acs.StartBit, c.MaskToValue(acs.StartBit, acs.EndBit),
			)
		default:
			panic("unimplemented")
		}
	}
	cFmts.Code += "\treturn 0;\n"
}

func genProcReturnsAccessBlockRead(p *fn.Proc, blk *fn.Block, cFmts *BlockCFormatters) {
	panic("unimplemented")
	/*
		cFmts.Code += fmt.Sprintf(
			"\t%s _rbuff[%d];\n", c.WidthToWriteType(blk.Width), p.ReturnsBufSize(),
		)
		cFmts.Code += fmt.Sprintf(
			"\tconst int err = iface.readb(%d, _rbuff, %d);\n", p.ReturnsStartAddr(), p.ReturnsBufSize(),
		)
		cFmts.Code += "\tif (err)\n\t\t return err;\n"
	*/
}
