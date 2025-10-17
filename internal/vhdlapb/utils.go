package vhdlapb

import (
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/types"
)

func addrRange(start, end int64, blk *fn.Block) types.SingleRange {
	return types.SingleRange{
		Start: start, End: end,
	}.Shift(-blk.StartAddr())
}
