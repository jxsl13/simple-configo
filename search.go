package configo

// binarySearchRangeInt requires a sorted list of ranges
func binarySearchRangeInt(a []intRange, x int) bool {
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
