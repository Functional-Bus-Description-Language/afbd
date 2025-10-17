package vhdlapb

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genMask(blk *fn.Block, mask *fn.Mask, fmts *BlockEntityFormatters) {
	if mask.IsArray {
		genMaskArray(blk, mask, fmts)
	} else {
		genMaskSingle(blk, mask, fmts)
	}
}

func genMaskArray(blk *fn.Block, mask *fn.Mask, fmts *BlockEntityFormatters) {
	panic("unimplemented")
}

func genMaskSingle(blk *fn.Block, mask *fn.Mask, fmts *BlockEntityFormatters) {
	dflt := ""
	if mask.InitValue != "" {
		dflt = fmt.Sprintf(" := %s", mask.InitValue.Extend(mask.Width))
	}

	fmts.EntityFunctionalPorts += fmt.Sprintf(
		";\n  %s_o : buffer std_logic_vector(%d downto 0)%s",
		mask.Name, mask.Width-1, dflt,
	)

	switch mask.Access.Type {
	case "SingleOneReg":
		genMaskSingleOneReg(blk, mask, fmts)
	case "SingleNRegs":
		genMaskSingleNRegs(blk, mask, fmts)
	default:
		panic("unimplemented")
	}
}

func genMaskSingleOneReg(blk *fn.Block, mask *fn.Mask, fmts *BlockEntityFormatters) {
	acs := mask.Access

	code := fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_o <= apb_req.wdata(%[2]d downto %[3]d);
    end if;
    apb_com.rdata(%[2]d downto %[3]d) <= %[1]s_o;`,
		mask.Name, acs.EndBit, acs.StartBit,
	)

	fmts.RegistersAccess.add(addrRange(acs.StartAddr, acs.EndAddr, blk), code)
}

func genMaskSingleNRegs(blk *fn.Block, mask *fn.Mask, fmts *BlockEntityFormatters) {
	if mask.Atomic {
		genMaskSingleNRegsAtomic(blk, mask, fmts)
	} else {
		genMaskSingleNRegsNonAtomic(blk, mask, fmts)
	}
}

func genMaskSingleNRegsAtomic(blk *fn.Block, mask *fn.Mask, fmts *BlockEntityFormatters) {
	acs := mask.Access
	strategy := SeparateLast
	atomicShadowRange := [2]int64{mask.Width - 1 - acs.EndRegWidth, 0}
	chunks := makeAccessChunksContinuous(acs, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		mask.Name, atomicShadowRange[0], atomicShadowRange[1],
	)

	for i, c := range chunks {
		var code string
		if (strategy == SeparateFirst && i == 0) || (strategy == SeparateLast && i == len(chunks)-1) {
			code = fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_o(%[2]s downto %[3]s) <= apb_req.wdata(%[4]d downto %[5]d);
      %[1]s_o(%[6]d downto %[7]d) <= %[1]s_atomic(%[6]d downto %[7]d);
    end if;
    apb_com.rdata(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);`,
				mask.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
				atomicShadowRange[0], atomicShadowRange[1],
			)
		} else {
			code = fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_atomic(%[2]s downto %[3]s) <= apb_req.wdata(%[4]d downto %[5]d);
    end if;
    apb_com.rdata(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);
`,
				mask.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
			)
		}

		fmts.RegistersAccess.add(c.addr.Shift(-blk.StartAddr()), code)
	}
}

func genMaskSingleNRegsNonAtomic(blk *fn.Block, mask *fn.Mask, fmts *BlockEntityFormatters) {
	acs := mask.Access
	chunks := makeAccessChunksContinuous(acs, Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_o(%[2]s downto %[3]s) <= apb_req.wdata(%[4]d downto %[5]d);
    end if;
    apb_com.rdata(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);`,
			mask.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
		)

		fmts.RegistersAccess.add(c.addr.Shift(-blk.StartAddr()), code)
	}
}
