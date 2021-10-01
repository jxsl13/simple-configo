package parsers

import (
	"sort"
	"strconv"
)

func intListToStringList(list []int) []string {
	result := make([]string, 0, len(list))
	for _, i := range list {
		result = append(result, strconv.Itoa(i))
	}
	return result
}

func floatListToStringList(list []float64, bitSize int) []string {
	result := make([]string, 0, len(list))
	for _, f := range list {
		result = append(result, strconv.FormatFloat(f, 'f', -1, bitSize))
	}
	return result
}

func listToSortedUniqueListString(list []string) []string {
	m := make(map[string]bool, len(list))
	for _, element := range list {
		m[element] = true
	}
	return setToSortedListString(m)
}

func listToSortedUniqueListInt(list []int) []int {
	m := make(map[int]bool, len(list))
	for _, element := range list {
		m[element] = true
	}
	return setToSortedListInt(m)
}

func listToSortedUniqueListFloat(list []float64) []float64 {
	m := make(map[float64]bool, len(list))
	for _, element := range list {
		m[element] = true
	}
	return setToSortedListFloat(m)
}

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
