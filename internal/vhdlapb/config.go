package vhdlapb

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
)

func genConfig(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	if cfg.IsArray {
		genConfigArray(blk, cfg, fmts)
	} else {
		genConfigSingle(blk, cfg, fmts)
	}
}

func genConfigArray(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	port := fmt.Sprintf(
		";\n  %s_o : buffer slv_vector(%d downto 0)(%d downto 0)",
		cfg.Name, cfg.Count-1, cfg.Width-1,
	)
	if cfg.InitValue != "" {
		port += fmt.Sprintf(" := (others => %s)", cfg.InitValue.Extend(cfg.Width))
	}
	fmts.EntityFunctionalPorts += port

	switch cfg.Access.Type {
	case "ArrayOneReg":
		genConfigArrayOneReg(blk, cfg, fmts)
	case "ArrayOneInReg":
		genConfigArrayOneInReg(blk, cfg, fmts)
	case "ArrayNInReg":
		genConfigArrayNInReg(blk, cfg, fmts)
	case "ArrayNInRegMInEndReg":
		genConfigArrayNInRegMInEndReg(blk, cfg, fmts)
	default:
		panic("unimplemented")
	}
}

func genConfigSingle(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	dflt := ""
	if cfg.InitValue != "" {
		dflt = fmt.Sprintf(" := %s", cfg.InitValue.Extend(cfg.Width))
	}

	fmts.EntityFunctionalPorts += fmt.Sprintf(
		";\n  %s_o : buffer std_logic_vector(%d downto 0)%s",
		cfg.Name, cfg.Width-1, dflt,
	)

	switch cfg.Access.Type {
	case "SingleOneReg":
		genConfigSingleOneReg(blk, cfg, fmts)
	case "SingleNRegs":
		genConfigSingleNRegs(blk, cfg, fmts)
	default:
		panic("unimplemented")
	}
}

func genConfigSingleOneReg(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access

	code := fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_o <= apb_req.wdata(%[2]d downto %[3]d);
    end if;
    apb_com.rdata(%[2]d downto %[3]d) <= %[1]s_o;`,
		cfg.Name, acs.EndBit, acs.StartBit,
	)

	fmts.RegistersAccess.add(addrRange(acs.StartAddr, acs.EndAddr, blk), code)
}

func genConfigSingleNRegs(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	if cfg.Atomic {
		genConfigSingleNRegsAtomic(blk, cfg, fmts)
	} else {
		genConfigSingleNRegsNonAtomic(blk, cfg, fmts)
	}
}

func genConfigSingleNRegsAtomic(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access
	strategy := SeparateLast
	atomicShadowRange := [2]int64{cfg.Width - 1 - acs.EndRegWidth, 0}
	chunks := makeAccessChunksContinuous(acs, strategy)

	fmts.SignalDeclarations += fmt.Sprintf(
		"signal %s_atomic : std_logic_vector(%d downto %d);\n",
		cfg.Name, atomicShadowRange[0], atomicShadowRange[1],
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
				cfg.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
				atomicShadowRange[0], atomicShadowRange[1],
			)
		} else {
			code = fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_atomic(%[2]s downto %[3]s) <= apb_req.wdata(%[4]d downto %[5]d);
    end if;
    apb_com.rdata(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);
`,
				cfg.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
			)
		}

		fmts.RegistersAccess.add(addrRange(c.addr[0], c.addr[1], blk), code)
	}
}

func genConfigSingleNRegsNonAtomic(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access
	chunks := makeAccessChunksContinuous(acs, Compact)

	for _, c := range chunks {
		code := fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_o(%[2]s downto %[3]s) <= apb_req.wdata(%[4]d downto %[5]d);
    end if;
    apb_com.rdata(%[4]d downto %[5]d) <= %[1]s_o(%[2]s downto %[3]s);`,
			cfg.Name, c.range_[0], c.range_[1], c.endBit, c.startBit,
		)

		fmts.RegistersAccess.add(addrRange(c.addr[0], c.addr[1], blk), code)
	}
}

func genConfigArrayOneInReg(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access

	code := fmt.Sprintf(`
    if apb_req.write = '1' then
      %[1]s_o(addr - %[2]d) <= apb_req.wdata(%[3]d downto %[4]d);
    end if;
    apb_com.rdata(%[3]d downto %[4]d) <= %[1]s_o(addr - %[2]d);`,
		cfg.Name, acs.StartAddr, acs.EndBit, acs.StartBit,
	)

	fmts.RegistersAccess.add(
		addrRange(acs.StartAddr, acs.StartAddr+acs.RegCount-1, blk),
		code,
	)
}

func genConfigArrayOneReg(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access

	code := fmt.Sprintf(`
    for i in 0 to %[1]d loop
      if apb_req.write = '1' then
        %[2]s_o(i) <= apb_req.wdata(%[3]d*(i+1)+%[4]d-1 downto %[3]d*i+%[4]d);
      end if;
      apb_com.rdata(%[3]d*(i+1)+%[4]d-1 downto %[3]d*i+%[4]d) <= %[2]s_o(i);
    end loop;`,
		cfg.Count-1, cfg.Name, acs.ItemWidth, acs.StartBit,
	)

	fmts.RegistersAccess.add(addrRange(acs.StartAddr, acs.EndAddr, blk), code)
}

func genConfigArrayNInReg(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access

	itemsInReg := acs.RegWidth / acs.ItemWidth

	code := fmt.Sprintf(`
    for i in 0 to %[1]d loop
      if apb_req.write = '1' then
        %[4]s_o((addr-%[5]d)*%[6]d+i) <= apb_req.wdata(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d);
      end if;
      apb_com.rdata(%[2]d*(i+1)+%[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_o((addr-%[5]d)*%[6]d+i);
    end loop;`,
		itemsInReg-1, acs.ItemWidth, acs.StartBit, cfg.Name, acs.StartAddr, itemsInReg,
	)

	fmts.RegistersAccess.add(addrRange(acs.StartAddr, acs.EndAddr, blk), code)
}

func genConfigArrayNInRegMInEndReg(blk *fn.Block, cfg *fn.Config, fmts *BlockEntityFormatters) {
	acs := cfg.Access

	itemsInReg := acs.RegWidth / acs.ItemWidth
	itemsInEndReg := acs.ItemCount - ((acs.RegCount - 1) * itemsInReg)

	code := fmt.Sprintf(`
    for i in 0 to %[1]d loop
      if apb_req.write = '1' then
        %[4]s_o((addr-%[5]d)*%[6]d+i) <= apb_req.wdata(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d);
      end if;
      apb_com.rdata(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i + %[3]d) <= %[4]s_o((addr-%[5]d)*%[6]d+i);
    end loop;`,
		itemsInReg-1, acs.ItemWidth, acs.StartBit, cfg.Name, acs.StartAddr, itemsInReg,
	)
	fmts.RegistersAccess.add(addrRange(acs.StartAddr, acs.EndAddr-1, blk), code)

	code = fmt.Sprintf(`
    for i in 0 to %[1]d loop
      if apb_req.write = '1' then
        %[4]s_o(%[5]d+i) <= apb_req.wdata(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d);
      end if;
      apb_com.rdata(%[2]d*(i+1) + %[3]d-1 downto %[2]d*i+%[3]d) <= %[4]s_o(%[5]d+i);
    end loop;`,
		itemsInEndReg-1, acs.ItemWidth, acs.StartBit, cfg.Name, (acs.RegCount-1)*itemsInReg,
	)

	fmts.RegistersAccess.add(addrRange(acs.EndAddr, acs.EndAddr, blk), code)
}
