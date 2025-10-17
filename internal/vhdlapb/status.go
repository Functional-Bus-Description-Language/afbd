package vhdlapb

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/types"
)

func genStatus(st *fn.Status, fmts *BlockEntityFormatters) {
	if st.IsArray {
		genStatusArray(st, fmts)
	} else {
		genStatusSingle(st, fmts)
	}
}

func genStatusArray(st *fn.Status, fmts *BlockEntityFormatters) {
	fmts.EntityFunctionalPorts += fmt.Sprintf(
		";\n  %s_i : in slv_vector(%d downto 0)(%d downto 0)",
		st.Name, st.Count-1, st.Width-1,
	)

	switch st.Access.Type {
	case "ArrayOneReg":
		genStatusArrayOneReg(st, fmts)
	case "ArrayOneInReg":
		genStatusArrayOneInReg(st, fmts)
	case "ArrayNInReg":
		genStatusArrayNInReg(st, fmts)
	case "ArrayNInRegMInEndReg":
		genStatusArrayNInRegMInEndReg(st, fmts)
	case "ArrayOneInNRegs":
		genStatusArrayOneInNRegs(st, fmts)
	default:
		panic("unimplemented")
	}
}

func genStatusSingle(st *fn.Status, fmts *BlockEntityFormatters) {
	fmts.EntityFunctionalPorts += fmt.Sprintf(
		";\n  %s_i : in std_logic_vector(%d downto 0)",
		st.Name, st.Width-1,
	)

	switch st.Access.Type {
	case "SingleOneReg":
		genStatusSingleOneReg(st, fmts)
	case "SingleNRegs":
		genStatusSingleNRegs(st, fmts)
	default:
		panic("unimplemented")
	}
}

func genStatusSingleOneReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access

	code := fmt.Sprintf(
		"    apb_com.rdata(%d downto %d) <= %s_i;\n",
		acs.EndBit, acs.StartBit, st.Name,
	)
	fmts.RegistersAccess.add(acs.AddrRange(), code)
}

func genStatusSingleNRegs(st *fn.Status, fmts *BlockEntityFormatters) {
	if st.Atomic {
		genStatusSingleNRegsAtomic(st, fmts)
	} else {
		genStatusSingleNRegsNonAtomic(st, fmts)
	}
}

func genStatusSingleNRegsAtomic(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access
	strategy := SeparateFirst
	atomicShadowRange := [2]int64{st.Width - 1, acs.StartRegWidth}
	chunks := makeAccessChunksContinuous(acs, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		st.Name, atomicShadowRange[0], atomicShadowRange[1],
	)

	for i, c := range chunks {
		var code string
		if (strategy == SeparateFirst && i == 0) || (strategy == SeparateLast && i == len(chunks)-1) {
			code = fmt.Sprintf(`
    %[1]s_atomic(%[2]d downto %[3]d) <= %[1]s_i(%[2]d downto %[3]d);
    apb_com.rdata(%[4]d downto %[5]d) <= %[1]s_i(%[6]s downto %[7]s);`,
				st.Name, atomicShadowRange[0], atomicShadowRange[1],
				c.endBit, c.startBit, c.range_[0], c.range_[1],
			)
		} else {
			code = fmt.Sprintf(
				"    apb_com.rdata(%d downto %d) <= %s_atomic(%s downto %s);",
				c.endBit, c.startBit, st.Name, c.range_[0], c.range_[1],
			)
		}

		fmts.RegistersAccess.add(c.addr, code)
	}
}

func genStatusSingleNRegsNonAtomic(st *fn.Status, fmts *BlockEntityFormatters) {
	chunks := makeAccessChunksContinuous(st.Access, Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(
			"    apb_com.rdata(%d downto %d) <= %s_i(%s downto %s);",
			c.endBit, c.startBit, st.Name, c.range_[0], c.range_[1],
		)

		fmts.RegistersAccess.add(c.addr, code)
	}
}

func genStatusArrayOneInReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access

	code := fmt.Sprintf(
		"    apb_com.rdata(%d downto %d) <= %s_i(addr - %d);",
		acs.EndBit, acs.StartBit, st.Name, acs.StartAddr,
	)

	fmts.RegistersAccess.add(
		types.SingleRange{Start: acs.StartAddr, End: acs.StartAddr + acs.RegCount - 1},
		code,
	)
}

func genStatusArrayOneReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access

	code := fmt.Sprintf(`
    for i in 0 to %[1]d loop
      apb_com.rdata(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(i);
    end loop;`,
		st.Count-1, acs.ItemWidth, acs.StartBit, st.Name,
	)

	fmts.RegistersAccess.add(acs.AddrRange(), code)
}

func genStatusArrayNInReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access

	itemsInReg := acs.RegWidth / acs.ItemWidth

	code := fmt.Sprintf(`
    for i in 0 to %[1]d loop
      apb_com.rdata(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
    end loop;`,
		itemsInReg-1, acs.ItemWidth, acs.StartBit, st.Name, acs.StartAddr, itemsInReg,
	)

	fmts.RegistersAccess.add(acs.AddrRange(), code)
}

func genStatusArrayNInRegMInEndReg(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access

	itemsInReg := acs.RegWidth / acs.ItemWidth
	itemsInEndReg := acs.ItemCount - ((acs.RegCount - 1) * itemsInReg)

	code := fmt.Sprintf(`
    for i in 0 to %[1]d loop
      apb_com.rdata(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d) <= %[4]s_i((addr-%[5]d)*%[6]d+i);
    end loop;`,
		itemsInReg-1, acs.ItemWidth, acs.StartBit, st.Name, acs.StartAddr, itemsInReg,
	)
	fmts.RegistersAccess.add(
		types.SingleRange{Start: acs.StartAddr, End: acs.EndAddr - 1}, code,
	)

	code = fmt.Sprintf(`
    for i in 0 to %[1]d loop
      apb_com.rdata(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_i(%[5]d+i);
    end loop;`,
		itemsInEndReg-1, acs.ItemWidth, acs.StartBit, st.Name, (acs.RegCount-1)*itemsInReg,
	)

	fmts.RegistersAccess.add(types.SingleRange{Start: acs.EndAddr, End: acs.EndAddr}, code)
}

func genStatusArrayOneInNRegs(st *fn.Status, fmts *BlockEntityFormatters) {
	if st.Atomic {
		panic("unimplemented")
	} else {
		genStatusArrayOneInNRegsNonAtomic(st, fmts)
	}
}

func genStatusArrayOneInNRegsNonAtomic(st *fn.Status, fmts *BlockEntityFormatters) {
	acs := st.Access

	regsPerItem := acs.ItemWidth / acs.RegWidth
	if acs.ItemWidth%acs.RegWidth != 0 {
		regsPerItem++
	}

	idx := fmt.Sprintf("(addr - %d) / %d", acs.StartAddr, regsPerItem)
	bite := fmt.Sprintf("(addr - %d) mod %d", acs.StartAddr, regsPerItem)
	lowerBound := fmt.Sprintf("(%s) * %d", bite, busWidth)
	upperBound := fmt.Sprintf("(%s) + %d", bite, busWidth-1)
	code := fmt.Sprintf(`
    if %[1]s = %[2]d then
      apb_com.rdata(%[3]d downto 0) <= %[4]s_i(%[5]s)(%[6]d downto %[7]s);
    else
      apb_com.rdata <= %[4]s_i(%[5]s)(%[8]s downto %[7]s);
    end if;`,
		bite, regsPerItem-1, acs.EndBit, st.Name, idx, st.Width-1, lowerBound, upperBound,
	)

	fmts.RegistersAccess.add(acs.AddrRange(), code)
}
