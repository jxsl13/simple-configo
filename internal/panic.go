package internal

// PanifIfNil is a check that enforces the user to pass a valid pointer value
func PanicIfNil(is ...interface{}) {
	for _, i := range is {
		if i == nil {
			panic("nil pointer parameter")
		}
	}
}

func PanicIfEmptyInt(list []int) {
	if len(list) == 0 {
		panic("list must not be empty")
	}
}

func PanicIfEmptyString(list []string) {
	if len(list) == 0 {
		panic("list must not be empty")
	}
}

func PanicIfEmptyFloat(list []float64) {
	if len(list) == 0 {
		panic("list must not be empty")
	}
}
