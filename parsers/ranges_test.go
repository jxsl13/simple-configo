package parsers_test

import (
	"testing"

	"github.com/jxsl13/simple-configo/parsers"
)

func Test_IntRange_Contains(t *testing.T) {
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
			ir := &parsers.IntRange{
				Min: tt.fields.Min,
				Max: tt.fields.Max,
			}
			if got := ir.Contains(tt.args.i); got != tt.want {
				t.Errorf("IntRange.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dDistinctRangeListInt_Contains(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name             string
		distictRangeList parsers.DistinctRangeListInt
		args             args
		want             bool
	}{
		{"#1", parsers.NewDistinctRangeListInt(1, 3), args{1}, true},
		{"#2", parsers.NewDistinctRangeListInt(1, 3), args{2}, true},
		{"#3", parsers.NewDistinctRangeListInt(1, 3), args{3}, true},
		{"#4", parsers.NewDistinctRangeListInt(1, 3), args{0}, false},
		{"#5", parsers.NewDistinctRangeListInt(1, 3), args{4}, false},
		{"#6", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{0}, false},
		{"#7", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{1}, true},
		{"#8", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{2}, true},
		{"#9", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{3}, true},
		{"#10", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{4}, true},
		{"#11", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{5}, true},
		{"#12", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{6}, true},
		{"#13", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{7}, true},
		{"#14", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{8}, true},
		{"#15", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{9}, true},
		{"#16", parsers.NewDistinctRangeListInt(1, 3, 2, 9), args{10}, false},

		{"#17", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{0}, false},
		{"#18", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{1}, true},
		{"#19", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{2}, true},
		{"#20", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{3}, true},
		{"#21", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{4}, true},
		{"#22", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{5}, true},
		{"#23", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{6}, true},
		{"#24", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{7}, true},
		{"#25", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{8}, true},
		{"#26", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{9}, true},
		{"#27", parsers.NewDistinctRangeListInt(3, 1, 9, 2), args{10}, false},

		{"#28", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{0}, false},
		{"#29", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{1}, true},
		{"#30", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{2}, true},
		{"#31", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{3}, true},
		{"#32", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{4}, true},
		{"#33", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{9}, true},
		{"#34", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{10}, true},
		{"#35", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{11}, true},
		{"#36", parsers.NewDistinctRangeListInt(3, 11, 1, 3, 2, 9), args{12}, false},

		{"#37", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{0}, false},
		{"#38", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{1}, true},
		{"#39", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{2}, true},
		{"#40", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{3}, true},
		{"#41", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{4}, true},
		{"#42", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{9}, true},
		{"#43", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{10}, true},
		{"#44", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{11}, true},
		{"#45", parsers.NewDistinctRangeListInt(10, 11, 1, 3, 2, 9), args{12}, false},

		{"#46", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{0}, false},
		{"#47", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{1}, true},
		{"#48", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{2}, true},
		{"#49", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{4}, true},
		{"#50", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{7}, true},
		{"#51", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{8}, true},
		{"#52", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{9}, false},
		{"#53", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{10}, false},
		{"#54", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{11}, true},
		{"#55", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{12}, true},
		{"#56", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{13}, true},
		{"#57", parsers.NewDistinctRangeListInt(11, 13, 1, 2, 2, 8), args{14}, false},

		{"#58", parsers.NewDistinctRangeListInt(-5, 3), args{-6}, false},
		{"#59", parsers.NewDistinctRangeListInt(-5, 3), args{-1}, true},
		{"#60", parsers.NewDistinctRangeListInt(-5, 3), args{-2}, true},
		{"#61", parsers.NewDistinctRangeListInt(-5, 3), args{-3}, true},
		{"#62", parsers.NewDistinctRangeListInt(-5, 3), args{0}, true},
		{"#63", parsers.NewDistinctRangeListInt(-5, 3), args{1}, true},
		{"#64", parsers.NewDistinctRangeListInt(-5, 3), args{2}, true},
		{"#65", parsers.NewDistinctRangeListInt(-5, 3), args{3}, true},
		{"#66", parsers.NewDistinctRangeListInt(-5, 3), args{4}, false},

		{"#67", parsers.NewDistinctRangeListInt(1, 5, 1, 3), args{0}, false},
		{"#68", parsers.NewDistinctRangeListInt(1, 5, 1, 3), args{1}, true},
		{"#69", parsers.NewDistinctRangeListInt(1, 5, 1, 3), args{2}, true},
		{"#70", parsers.NewDistinctRangeListInt(1, 5, 1, 3), args{3}, true},
		{"#71", parsers.NewDistinctRangeListInt(1, 5, 1, 3), args{4}, true},
		{"#72", parsers.NewDistinctRangeListInt(1, 5, 1, 3), args{5}, true},
		{"#73", parsers.NewDistinctRangeListInt(1, 5, 1, 3), args{6}, false},

		{"#74", parsers.NewDistinctRangeListInt(1, 5, 2, 3), args{6}, false},
		{"#75", parsers.NewDistinctRangeListInt(1, 5, 2, 3), args{2}, true},
		{"#76", parsers.NewDistinctRangeListInt(1, 5, 2, 3), args{4}, true},

		//{"#99", parsers.NewDistinctRangeListInt(1, 99), args{2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &tt.distictRangeList
			t.Logf("LIST: %v", d)
			if got := d.Contains(tt.args.i); got != tt.want {
				t.Errorf("DistinctRangeListInt.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
