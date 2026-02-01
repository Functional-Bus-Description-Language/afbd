package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/afbd/internal/c"
	"github.com/Functional-Bus-Description-Language/afbd/internal/utils"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
)

func genConfig(cfg *fn.Config, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	if cfg.IsArray {
		panic("unimplemented")
	} else {
		genConfigSingle(cfg, blk, hFmts, cFmts)
	}
}

func genConfigSingle(cfg *fn.Config, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	switch cfg.Access.Type {
	case "SingleOneReg":
		genConfigSingleOneReg(cfg, blk, hFmts, cFmts)
	default:
		panic("unimplemented")
	}
}

func genConfigSingleOneReg(cfg *fn.Config, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	rType := c.WidthToReadType(cfg.Width)
	wType := c.WidthToWriteType(cfg.Width)

	readSignature := fmt.Sprintf(
		"int afbd_%s_%s_read(afbd_iface_t * const iface, %s const data)",
		hFmts.BlockName, cfg.Name, rType.String(),
	)
	writeSignature := fmt.Sprintf(
		"int afbd_%s_%s_write(afbd_iface_t * const iface, %s const data)",
		hFmts.BlockName, cfg.Name, wType.String(),
	)

	hFmts.Code += fmt.Sprintf("\n%s;\n%s;\n", readSignature, writeSignature)

	acs := cfg.Access
	addr := acs.StartAddr
	if !args.CSync.OffsetAddr {
		addr += blk.StartAddr()
	}

	cFmts.Code += fmt.Sprintf("\n%s\n{\n", readSignature)
	if readType.Typ() != "ByteArray" && rType.Typ() != "ByteArray" {
		if busWidth == cfg.Width {
			cFmts.Code += fmt.Sprintf(
				"\treturn iface->read(iface, %d, data);\n};\n", addr,
			)
		} else {
			cFmts.Code += fmt.Sprintf(`	%s aux;
	const int err = iface->read(iface, %d, &aux);
	if (err)
		return err;
	*data = (aux >> %d) & 0x%x;
	return 0;
};
`, readType.Depointer().String(), addr, acs.StartBit, utils.Uint64Mask(acs.StartBit, acs.EndBit),
			)
		}
	} else {
		panic("unimplemented")
	}

	cFmts.Code += fmt.Sprintf("\n%s\n{\n", writeSignature)
	if readType.Typ() != "ByteArray" && rType.Typ() != "ByteArray" {
		if busWidth == cfg.Width {
			cFmts.Code += fmt.Sprintf(
				"\treturn iface->write(iface, %d, data);\n};\n", addr,
			)
		} else {
			cFmts.Code += fmt.Sprintf(
				"	return iface->write(iface, %d, (data << %d));\n};", addr, acs.StartBit,
			)
		}
	} else {
		panic("unimplemented")
	}
}
