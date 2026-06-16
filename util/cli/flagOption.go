package cli

import (
	"net"
	"time"

	"github.com/spf13/pflag"
)

type FlagVar interface {
	bool | []bool |
		string | []string | map[string]string |
		int | int8 | int16 | int32 | int64 | []int | []int32 | []int64 |
		uint | uint8 | uint16 | uint32 | uint64 | []uint |
		float32 | float64 | []float32 | []float64 |
		net.IP | []net.IP |
		time.Duration | []time.Duration
}

type FlagOption func(*pflag.FlagSet)

func WithFlagVar[T FlagVar](p *T, name string, defVal T, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		switch v := any(defVal).(type) {
		case bool:
			f.BoolVar(any(p).(*bool), name, v, usage)

		case []bool:
			f.BoolSliceVar(any(p).(*[]bool), name, v, usage)

		case string:
			f.StringVar(any(p).(*string), name, v, usage)

		case []string:
			f.StringSliceVar(any(p).(*[]string), name, v, usage)

		case map[string]string:
			f.StringToStringVar(any(p).(*map[string]string), name, v, usage)

		case int:
			f.IntVar(any(p).(*int), name, v, usage)

		case int8:
			f.Int8Var(any(p).(*int8), name, v, usage)

		case int16:
			f.Int16Var(any(p).(*int16), name, v, usage)

		case int32:
			f.Int32Var(any(p).(*int32), name, v, usage)

		case int64:
			f.Int64Var(any(p).(*int64), name, v, usage)

		case []int:
			f.IntSliceVar(any(p).(*[]int), name, v, usage)

		case []int32:
			f.Int32SliceVar(any(p).(*[]int32), name, v, usage)

		case []int64:
			f.Int64SliceVar(any(p).(*[]int64), name, v, usage)

		case uint:
			f.UintVar(any(p).(*uint), name, v, usage)

		case uint8:
			f.Uint8Var(any(p).(*uint8), name, v, usage)

		case uint16:
			f.Uint16Var(any(p).(*uint16), name, v, usage)

		case uint32:
			f.Uint32Var(any(p).(*uint32), name, v, usage)

		case uint64:
			f.Uint64Var(any(p).(*uint64), name, v, usage)

		case []uint:
			f.UintSliceVar(any(p).(*[]uint), name, v, usage)

		case float32:
			f.Float32Var(any(p).(*float32), name, v, usage)

		case float64:
			f.Float64Var(any(p).(*float64), name, v, usage)

		case []float32:
			f.Float32SliceVar(any(p).(*[]float32), name, v, usage)

		case []float64:
			f.Float64SliceVar(any(p).(*[]float64), name, v, usage)

		case net.IP:
			f.IPVar(any(p).(*net.IP), name, v, usage)

		case []net.IP:
			f.IPSliceVar(any(p).(*[]net.IP), name, v, usage)

		case time.Duration:
			f.DurationVar(any(p).(*time.Duration), name, v, usage)
		}
	}
}

func WithFlagVarP[T FlagVar](p *T, name string, short string, value T, usage string) FlagOption {
	return func(f *pflag.FlagSet) {
		switch v := any(value).(type) {
		case bool:
			f.BoolVarP(any(p).(*bool), name, short, v, usage)

		case []bool:
			f.BoolSliceVarP(any(p).(*[]bool), name, short, v, usage)

		case string:
			f.StringVarP(any(p).(*string), name, short, v, usage)

		case []string:
			f.StringSliceVarP(any(p).(*[]string), name, short, v, usage)

		case map[string]string:
			f.StringToStringVarP(any(p).(*map[string]string), name, short, v, usage)

		case int:
			f.IntVarP(any(p).(*int), name, short, v, usage)

		case int8:
			f.Int8VarP(any(p).(*int8), name, short, v, usage)

		case int16:
			f.Int16VarP(any(p).(*int16), name, short, v, usage)

		case int32:
			f.Int32VarP(any(p).(*int32), name, short, v, usage)

		case int64:
			f.Int64VarP(any(p).(*int64), name, short, v, usage)

		case []int:
			f.IntSliceVarP(any(p).(*[]int), name, short, v, usage)

		case []int32:
			f.Int32SliceVarP(any(p).(*[]int32), name, short, v, usage)

		case []int64:
			f.Int64SliceVarP(any(p).(*[]int64), name, short, v, usage)

		case uint:
			f.UintVarP(any(p).(*uint), name, short, v, usage)

		case uint8:
			f.Uint8VarP(any(p).(*uint8), name, short, v, usage)

		case uint16:
			f.Uint16VarP(any(p).(*uint16), name, short, v, usage)

		case uint32:
			f.Uint32VarP(any(p).(*uint32), name, short, v, usage)

		case uint64:
			f.Uint64VarP(any(p).(*uint64), name, short, v, usage)

		case []uint:
			f.UintSliceVarP(any(p).(*[]uint), name, short, v, usage)

		case float32:
			f.Float32VarP(any(p).(*float32), name, short, v, usage)

		case float64:
			f.Float64VarP(any(p).(*float64), name, short, v, usage)

		case []float32:
			f.Float32SliceVarP(any(p).(*[]float32), name, short, v, usage)

		case []float64:
			f.Float64SliceVarP(any(p).(*[]float64), name, short, v, usage)

		case net.IP:
			f.IPVarP(any(p).(*net.IP), name, short, v, usage)

		case []net.IP:
			f.IPSliceVarP(any(p).(*[]net.IP), name, short, v, usage)

		case time.Duration:
			f.DurationVarP(any(p).(*time.Duration), name, short, v, usage)
		}
	}
}
