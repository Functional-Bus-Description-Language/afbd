package vhdlapb

import (
	"fmt"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/types"
)

type RegisterMap map[types.SingleRange]string

// addrRng is address range
func (regMap RegisterMap) add(addrRng types.SingleRange, code string) {
	if addrRng.End < addrRng.Start {
		panic(fmt.Sprintf("adrrRng.End (%d) < addrRng.Start (%d)", addrRng.End, addrRng.Start))
	}

	overlaps := []types.SingleRange{}
	for ar := range regMap {
		if (ar.Start <= addrRng.Start && addrRng.Start <= ar.End) ||
			ar.Start <= addrRng.End && addrRng.End <= ar.End {
			overlaps = append(overlaps, ar)
		}
	}

	if len(overlaps) == 0 {
		regMap[addrRng] = code
		return
	}

	if len(overlaps) == 1 && overlaps[0].Start == addrRng.Start && overlaps[0].End == addrRng.End {
		regMap[addrRng] += code
		return
	}

	for _, ovlp := range overlaps {
		ovlpCode := regMap[ovlp]
		delete(regMap, ovlp)

		// Middle overlap
		if ovlp.Start < addrRng.Start && addrRng.End < ovlp.End {
			regMap[types.SingleRange{
				Start: ovlp.Start, End: addrRng.Start - 1,
			}] = ovlpCode
			regMap[addrRng] = ovlpCode + code
			regMap[types.SingleRange{
				Start: addrRng.End + 1, End: ovlp.End,
			}] = ovlpCode
		}

		// Start overlap
		if addrRng.Start <= ovlp.Start && addrRng.End < ovlp.End {
			regMap[types.SingleRange{
				Start: addrRng.End + 1, End: ovlp.End,
			}] = ovlpCode

			if ovlp.Start == addrRng.Start {
				regMap[addrRng] = ovlpCode + code
			} else {
				regMap[types.SingleRange{
					Start: addrRng.Start, End: ovlp.Start - 1,
				}] = code
				regMap[types.SingleRange{
					Start: ovlp.Start, End: addrRng.End,
				}] = ovlpCode + code
			}
		}

		// End overlap
		if ovlp.Start < addrRng.Start && ovlp.End <= addrRng.End {
			regMap[types.SingleRange{
				Start: ovlp.Start, End: addrRng.Start - 1,
			}] = ovlpCode

			if ovlp.End == addrRng.End {
				regMap[addrRng] = ovlpCode + code
			} else {
				regMap[types.SingleRange{
					Start: addrRng.Start, End: ovlp.End,
				}] = ovlpCode + code
				regMap[types.SingleRange{
					Start: ovlp.End + 1, End: addrRng.End,
				}] = code
			}
		}
	}
}
