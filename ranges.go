package configo

import (
	"fmt"
	"sort"
	"strings"
)

type intRange struct {
	Min int
	Max int
}

func (ir *intRange) String() string {
	return fmt.Sprintf("[%d:%d]", ir.Min, ir.Max)
}

func (ir *intRange) Contains(i int) bool {
	return ir.Min <= i && i <= ir.Max
}

func (ir *intRange) Below(i int) bool {
	return ir.Max < i
}

func (ir *intRange) Above(i int) bool {
	return i < ir.Min
}

// sorting
type byIntRange []intRange

func (a byIntRange) Len() int           { return len(a) }
func (a byIntRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIntRange) Less(i, j int) bool { return a[i].Min < a[j].Min }

type distinctRangeListInt struct {
	r        []intRange
	isSorted bool
}

func (d *distinctRangeListInt) Contains(i int) bool {
	if d.isSorted {
		return binarySearchRangeInt(d.r, i)
	}

	sort.Sort(byIntRange(d.r))
	d.isSorted = true
	return binarySearchRangeInt(d.r, i)
}

func (d *distinctRangeListInt) String() string {
	var sb strings.Builder

	for idx, r := range d.r {
		sb.WriteString(r.String())
		if idx < len(d.r)-1 {
			sb.WriteString(", ")
		}
	}

}

func newDistinctRangeListInt(minMaxRanges ...int) distinctRangeListInt {
	if len(minMaxRanges)%2 != 0 {
		panic(fmt.Errorf("passed parameter list 'minMaxRanges' must contain an even number of parameters"))
	}

	rangesList := make([]intRange, 0, len(minMaxRanges)/2)

	for i := 0; i < len(minMaxRanges); i += 2 {
		min := minMaxRanges[i]
		max := minMaxRanges[i+1]
		if max < min {
			min, max = max, min
		}
		rangesList = append(rangesList, intRange{min, max})
	}

	distinctList := make([]intRange, 0, len(rangesList)/2)
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
		// TODO: testing & c&p for floatRage
		distinctList = append(distinctList, *currentRange)
	}

	return distinctRangeListInt{distinctList, true}
}

// TODO:
type floatRange struct {
	Min float64
	Max float64
}

func (fr *floatRange) Contains(f float64) bool {
	return fr.Min <= f && f <= fr.Max
}

func (fr *floatRange) Below(f float64) bool {
	return fr.Max < f
}

func (fr *floatRange) Above(f float64) bool {
	return f < fr.Min
}

// sorting
type byFloatRange []intRange

func (a byFloatRange) Len() int           { return len(a) }
func (a byFloatRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFloatRange) Less(i, j int) bool { return a[i].Max < a[j].Max }
