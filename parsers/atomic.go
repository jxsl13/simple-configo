package parsers

import (
	"math"
	"strconv"
	"sync/atomic"
	"unsafe"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// AtomicInt32 sets the out variable atomically.
func AtomicInt32(out *int32) configo.ParserFunc {
	internal.PanicIfNil(out)
	return func(value string) error {
		i32, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		atomic.StoreInt32(out, int32(i32))
		return nil
	}
}

// AtomicInt64 sets the out variable atomically.
func AtomicInt64(out *int64) configo.ParserFunc {
	internal.PanicIfNil(out)
	return func(value string) error {
		i64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		atomic.StoreInt64(out, i64)
		return nil
	}
}

// AtomicFloat32 sets the out variable atomically.
func AtomicFloat32(out *float32) configo.ParserFunc {
	internal.PanicIfNil(out)
	return func(value string) error {
		f32, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}

		u32 := (*uint32)(unsafe.Pointer(out))
		atomic.StoreUint32(u32, math.Float32bits(float32(f32)))
		return nil
	}
}

// AtomicFloat64 sets the out variable atomically.
func AtomicFloat64(out *float64) configo.ParserFunc {
	internal.PanicIfNil(out)
	return func(value string) error {
		f64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		u64 := (*uint64)(unsafe.Pointer(out))
		atomic.StoreUint64(u64, math.Float64bits(f64))
		return nil
	}
}
