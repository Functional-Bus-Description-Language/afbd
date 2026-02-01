package csync

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/afbd/internal/c"
	"github.com/Functional-Bus-Description-Language/afbd/internal/utils"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
)

func genStatus(st *fn.Status, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	if st.IsArray {
		panic("unimplemented")
	} else {
		genStatusSingle(st, blk, hFmts, cFmts)
	}
}

func genStatusSingle(st *fn.Status, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	switch st.Access.Type {
	case "SingleOneReg":
		genStatusSingleOneReg(st, blk, hFmts, cFmts)
	default:
		panic("unimplemented")
	}
}

func genStatusSingleOneReg(st *fn.Status, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	typ := c.WidthToReadType(st.Width)
	signature := fmt.Sprintf(
		"int afbd_%s_%s_read(afbd_iface_t * const iface, %s const data)",
		hFmts.BlockName, st.Name, typ.String(),
	)

	hFmts.Code += fmt.Sprintf("\n%s;\n", signature)

	acs := st.Access
	addr := acs.StartAddr
	if !args.CSync.OffsetAddr {
		addr += blk.StartAddr()
	}

	cFmts.Code += fmt.Sprintf("\n%s\n{\n", signature)
	if readType.Typ() != "ByteArray" && typ.Typ() != "ByteArray" {
		if busWidth == st.Width {
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
}
