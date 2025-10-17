package vhdlapb

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/types"
)

type chunkStrategy uint8

const (
	Compact chunkStrategy = iota // Use only for non atomic elements.
	SeparateFirst
	SeparateLast
)

type accessChunk struct {
	addr     types.SingleRange
	range_   [2]string
	startBit int64
	endBit   int64
}

func makeAccessChunksContinuous(acs types.Access, strategy chunkStrategy) []accessChunk {
	startBit := acs.StartBit
	endBit := acs.EndBit

	cs := []accessChunk{}

	if strategy == Compact && acs.StartRegWidth == busWidth && acs.EndRegWidth == busWidth {
		cs = append(cs, accessChunk{
			addr: types.SingleRange{Start: acs.StartAddr, End: acs.EndAddr},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, acs.StartAddr),
				fmt.Sprintf("%d * (addr - %d)", busWidth, acs.StartAddr),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
	} else if acs.RegCount == 2 {
		cs = append(cs, accessChunk{
			addr:     types.SingleRange{Start: acs.StartAddr, End: acs.StartAddr},
			range_:   [2]string{fmt.Sprintf("%d", acs.StartRegWidth-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: types.SingleRange{Start: acs.EndAddr, End: acs.EndAddr},
			range_: [2]string{
				fmt.Sprintf("%d", acs.ItemWidth-1),
				fmt.Sprintf("%d", acs.ItemWidth-acs.EndRegWidth),
			},
			startBit: 0,
			endBit:   endBit,
		})
	} else if strategy == SeparateLast && acs.StartRegWidth == busWidth {
		cs = append(cs, accessChunk{
			addr: types.SingleRange{Start: acs.StartAddr, End: acs.EndAddr - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) - 1", busWidth, acs.StartAddr),
				fmt.Sprintf("%d * (addr - %d)", busWidth, acs.StartAddr),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: types.SingleRange{Start: acs.EndAddr, End: acs.EndAddr},
			range_: [2]string{
				fmt.Sprintf("%d", acs.ItemWidth-1),
				fmt.Sprintf("%d", acs.ItemWidth-acs.EndRegWidth),
			},
			startBit: 0,
			endBit:   endBit,
		})
	} else if strategy == SeparateFirst && acs.EndRegWidth == busWidth {
		cs = append(cs, accessChunk{
			addr:     types.SingleRange{Start: acs.StartAddr, End: acs.StartAddr},
			range_:   [2]string{fmt.Sprintf("%d", acs.StartRegWidth-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: types.SingleRange{Start: acs.StartAddr + 1, End: acs.EndAddr},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d + 1) + %d", busWidth, acs.StartAddr, acs.StartRegWidth-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, acs.StartAddr, acs.StartRegWidth),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
	} else {
		cs = append(cs, accessChunk{
			addr:     types.SingleRange{Start: acs.StartAddr, End: acs.StartAddr},
			range_:   [2]string{fmt.Sprintf("%d", acs.StartRegWidth-1), "0"},
			startBit: startBit,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: types.SingleRange{Start: acs.StartAddr + 1, End: acs.EndAddr - 1},
			range_: [2]string{
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, acs.StartAddr, acs.StartRegWidth-1),
				fmt.Sprintf("%d * (addr - %d) + %d", busWidth, acs.StartAddr+1, acs.StartRegWidth),
			},
			startBit: 0,
			endBit:   busWidth - 1,
		})
		cs = append(cs, accessChunk{
			addr: types.SingleRange{Start: acs.EndAddr, End: acs.EndAddr},
			range_: [2]string{
				fmt.Sprintf("%d", acs.ItemWidth-1),
				fmt.Sprintf("%d", acs.ItemWidth-acs.EndRegWidth),
			},
			startBit: 0,
			endBit:   endBit,
		})
	}

	return cs
}
