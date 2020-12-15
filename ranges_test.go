package configo

import (
	"testing"
)

func Test_intRange_Contains(t *testing.T) {
	type fields struct {
		Min int
		Max int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"#1", fields{1, 2}, args{1}, true},
		{"#2", fields{1, 2}, args{2}, true},
		{"#3", fields{1, 2}, args{0}, false},
		{"#4", fields{1, 2}, args{3}, false},
		{"#5", fields{1, 2}, args{-900}, false},
		{"#6", fields{1, 2}, args{+999999999999}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := &intRange{
				Min: tt.fields.Min,
				Max: tt.fields.Max,
			}
			if got := ir.Contains(tt.args.i); got != tt.want {
				t.Errorf("intRange.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_distinctRangeListInt_Contains(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name             string
		distictRangeList distinctRangeListInt
		args             args
		want             bool
	}{
		{"#1", newDistinctRangeListInt(1, 3), args{1}, true},
		{"#2", newDistinctRangeListInt(1, 3), args{2}, true},
		{"#3", newDistinctRangeListInt(1, 3), args{3}, true},
		{"#4", newDistinctRangeListInt(1, 3), args{0}, false},
		{"#5", newDistinctRangeListInt(1, 3), args{4}, false},
		{"#6", newDistinctRangeListInt(1, 3, 2, 9), args{0}, false},
		{"#7", newDistinctRangeListInt(1, 3, 2, 9), args{1}, true},
		{"#8", newDistinctRangeListInt(1, 3, 2, 9), args{2}, true},
		{"#9", newDistinctRangeListInt(1, 3, 2, 9), args{3}, true},
		{"#10", newDistinctRangeListInt(1, 3, 2, 9), args{4}, true},
		{"#11", newDistinctRangeListInt(1, 3, 2, 9), args{5}, true},
		{"#12", newDistinctRangeListInt(1, 3, 2, 9), args{6}, true},
		{"#13", newDistinctRangeListInt(1, 3, 2, 9), args{7}, true},
		{"#14", newDistinctRangeListInt(1, 3, 2, 9), args{8}, true},
		{"#15", newDistinctRangeListInt(1, 3, 2, 9), args{9}, true},
		{"#16", newDistinctRangeListInt(1, 3, 2, 9), args{10}, false},

		{"#17", newDistinctRangeListInt(3, 1, 9, 2), args{0}, false},
		{"#18", newDistinctRangeListInt(3, 1, 9, 2), args{1}, true},
		{"#19", newDistinctRangeListInt(3, 1, 9, 2), args{2}, true},
		{"#20", newDistinctRangeListInt(3, 1, 9, 2), args{3}, true},
		{"#21", newDistinctRangeListInt(3, 1, 9, 2), args{4}, true},
		{"#22", newDistinctRangeListInt(3, 1, 9, 2), args{5}, true},
		{"#23", newDistinctRangeListInt(3, 1, 9, 2), args{6}, true},
		{"#24", newDistinctRangeListInt(3, 1, 9, 2), args{7}, true},
		{"#25", newDistinctRangeListInt(3, 1, 9, 2), args{8}, true},
		{"#26", newDistinctRangeListInt(3, 1, 9, 2), args{9}, true},
		{"#27", newDistinctRangeListInt(3, 1, 9, 2), args{10}, false},

		{"#28", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{0}, false},
		{"#29", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{1}, true},
		{"#30", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{2}, true},
		{"#31", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{3}, true},
		{"#32", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{4}, true},
		{"#33", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{9}, true},
		{"#34", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{10}, true},
		{"#35", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{11}, true},
		{"#36", newDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{12}, false},

		{"#37", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{0}, false},
		{"#38", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{1}, true},
		{"#39", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{2}, true},
		{"#40", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{3}, true},
		{"#41", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{4}, true},
		{"#42", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{9}, true},
		{"#43", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{10}, true},
		{"#44", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{11}, true},
		{"#45", newDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{12}, false},

		{"#46", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{0}, false},
		{"#47", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{1}, true},
		{"#48", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{2}, true},
		{"#49", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{4}, true},
		{"#50", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{7}, true},
		{"#51", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{8}, true},
		{"#52", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{9}, false},
		{"#53", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{10}, false},
		{"#54", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{11}, true},
		{"#55", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{12}, true},
		{"#56", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{13}, true},
		{"#57", newDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{14}, false},

		{"#58", newDistinctRangeListInt(-5, 3), args{-6}, false},
		{"#59", newDistinctRangeListInt(-5, 3), args{-1}, true},
		{"#60", newDistinctRangeListInt(-5, 3), args{-2}, true},
		{"#61", newDistinctRangeListInt(-5, 3), args{-3}, true},
		{"#62", newDistinctRangeListInt(-5, 3), args{0}, true},
		{"#63", newDistinctRangeListInt(-5, 3), args{1}, true},
		{"#64", newDistinctRangeListInt(-5, 3), args{2}, true},
		{"#65", newDistinctRangeListInt(-5, 3), args{3}, true},
		{"#66", newDistinctRangeListInt(-5, 3), args{4}, false},

		{"#67", newDistinctRangeListInt(1, 5, 1, 3), args{0}, false},
		{"#68", newDistinctRangeListInt(1, 5, 1, 3), args{1}, true},
		{"#69", newDistinctRangeListInt(1, 5, 1, 3), args{2}, true},
		{"#70", newDistinctRangeListInt(1, 5, 1, 3), args{3}, true},
		{"#71", newDistinctRangeListInt(1, 5, 1, 3), args{4}, true},
		{"#72", newDistinctRangeListInt(1, 5, 1, 3), args{5}, true},
		{"#73", newDistinctRangeListInt(1, 5, 1, 3), args{6}, false},

		{"#74", newDistinctRangeListInt(1, 5, 2, 3), args{6}, false},
		{"#75", newDistinctRangeListInt(1, 5, 2, 3), args{2}, true},
		{"#76", newDistinctRangeListInt(1, 5, 2, 3), args{4}, true},

		//{"#99", newDistinctRangeListInt(1, 99), args{2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &tt.distictRangeList
			t.Logf("LIST: %v", d)
			if got := d.Contains(tt.args.i); got != tt.want {
				t.Errorf("distinctRangeListInt.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
