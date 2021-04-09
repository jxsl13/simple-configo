package parsers

import (
	"fmt"
	"sort"
	"strings"
)

// IntRange is a representation of a range that has a lower and upper bound (Min, Max).
// Min and Max are elements of the range.
type IntRange struct {
	Min int
	Max int
}

func (ir *IntRange) String() string {
	return fmt.Sprintf("[%d:%d]", ir.Min, ir.Max)
}

func (ir *IntRange) Contains(i int) bool {
	return ir.Min <= i && i <= ir.Max
}

func (ir *IntRange) Below(i int) bool {
	return ir.Max < i
}

func (ir *IntRange) Above(i int) bool {
	return i < ir.Min
}

// sorting
type byIntRange []IntRange

func (a byIntRange) Len() int      { return len(a) }
func (a byIntRange) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byIntRange) Less(i, j int) bool {
	if a[i].Min == a[j].Min {
		return a[i].Max < a[j].Max
	}
	return a[i].Min < a[j].Min
}

type DistinctRangeListInt struct {
	r []IntRange
}

func (d *DistinctRangeListInt) Contains(i int) bool {
	return binarySearchRangeInt(d.r, i)
}

func (d *DistinctRangeListInt) String() string {
	var sb strings.Builder
	const expectedChars = 7
	sb.Grow(expectedChars * len(d.r))

	for idx, r := range d.r {
		sb.WriteString(r.String())
		if idx < len(d.r)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func NewDistinctRangeListInt(minMaxRanges ...int) DistinctRangeListInt {
	if len(minMaxRanges)%2 != 0 {
		panic(fmt.Errorf("passed parameter list 'minMaxRanges' must contain an even number of parameters"))
	}

	rangesList := make([]IntRange, 0, len(minMaxRanges)/2)

	for i := 0; i < len(minMaxRanges); i += 2 {
		min := minMaxRanges[i]
		max := minMaxRanges[i+1]
		if max < min {
			min, max = max, min
		}
		rangesList = append(rangesList, IntRange{Min: min, Max: max})
	}

	distinctList := make([]IntRange, 0, len(rangesList)/2)
	sort.Sort(byIntRange(rangesList))
	for idx := range rangesList {
		currentRange := &rangesList[idx]
		if idx == 0 {
			distinctList = append(distinctList, *currentRange)
			continue
		}
		lastProcessedRange := &distinctList[len(distinctList)-1]

		// [1:12], [9:13]
		// [1:12], [2:12]
		if lastProcessedRange.Max >= currentRange.Min {
			if lastProcessedRange.Max > currentRange.Max {
				// skip, as element lies within previous range
				continue
			}
			// expand previous range to a lager range than before
			lastProcessedRange.Max = currentRange.Max
			// skip current range after updating previously processed one
			continue
		}
		distinctList = append(distinctList, *currentRange)
	}

	return DistinctRangeListInt{distinctList}
}

// binarySearchRangeInt requires a sorted list of ranges
func binarySearchRangeInt(a []IntRange, x int) bool {
	start := 0
	end := len(a) - 1
	for start <= end {
		mid := (start + end) / 2
		if a[mid].Contains(x) {
			return true
		} else if a[mid].Below(x) {
			start = mid + 1
		} else if a[mid].Above(x) {
			end = mid - 1
		}
	}
	return false
}
