package csync

import (
	"fmt"
	"strconv"

	"github.com/Functional-Bus-Description-Language/afbd/internal/c"
	"github.com/Functional-Bus-Description-Language/afbd/internal/utils"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/access"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genStatic(st *fn.Static, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	if st.IsArray {
		panic("unimplemented")
	} else {
		genStaticSingle(st, blk, hFmts, cFmts)
	}
}

func genStaticSingle(st *fn.Static, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	switch st.Access.(type) {
	case access.SingleOneReg:
		genStaticSingleOneReg(st, blk, hFmts, cFmts)
	default:
		panic("unimplemented")
	}
}

func genStaticSingleOneReg(st *fn.Static, blk *fn.Block, hFmts *BlockHFormatters, cFmts *BlockCFormatters) {
	wTyp := c.WidthToWriteType(st.Width)
	rTyp := c.WidthToReadType(st.Width)

	hFmts.Code += fmt.Sprintf(
		"\nextern const %s afbd_%s_%s;\n",
		wTyp.String(), hFmts.BlockName, st.Name,
	)

	signature := fmt.Sprintf(
		"int afbd_%s_%s_read(const afbd_iface_t * const iface, %s const data)",
		hFmts.BlockName, st.Name, rTyp.String(),
	)

	hFmts.Code += fmt.Sprintf("%s;\n", signature)

	cFmts.Code += fmt.Sprintf(
		"\nconst %s afbd_%s_%s = %s;\n",
		wTyp.String(), hFmts.BlockName, st.Name,
		// XXX: Uint64 is currently used. Below code needs fix if static is longer than 64 bits.
		fmt.Sprintf("0x%s", strconv.FormatUint(st.InitValue.Uint64(), 16)),
	)

	acs := st.Access.(access.SingleOneReg)
	cFmts.Code += fmt.Sprintf("%s {\n", signature)
	if readType.Typ() != "ByteArray" && rTyp.Typ() != "ByteArray" {
		if busWidth == st.Width {
			cFmts.Code += fmt.Sprintf(
				"\treturn iface->read(%d, data);\n};\n", blk.StartAddr()+acs.StartAddr(),
			)
		} else {
			cFmts.Code += fmt.Sprintf(`	%s aux;
	const int err = iface->read(%d, &aux);
	if (err)
		return err;
	*data = (aux >> %d) & 0x%x;
	return 0;
};
`, readType.Depointer().String(), blk.StartAddr()+acs.StartAddr(), acs.StartBit(), utils.Uint64Mask(acs.StartBit(), acs.EndBit()),
			)
		}
	} else {
		panic("unimplemented")
	}
}
