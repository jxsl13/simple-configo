package configo

import "sort"

func setToSortedListInt(a map[int]bool) (result []int) {
	result = make([]int, 0, len(a))
	for k := range a {
		result = append(result, k)
	}
	sort.Ints(result)
	return
}

func setToSortedListFloat(a map[float64]bool) (result []float64) {
	result = make([]float64, 0, len(a))
	for k := range a {
		result = append(result, k)
	}
	sort.Float64s(result)
	return
}

func setToSortedListString(a map[string]bool) (result []string) {
	result = make([]string, 0, len(a))
	for k := range a {
		result = append(result, k)
	}
	sort.Strings(result)
	return
}
