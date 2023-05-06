package stats

import (
	"encoding/json"
	"math"
	"strconv"
)

type StatName string
type Number interface {
	uint64 | float64
}

const (
	Count             StatName = "count"
	Min               StatName = "min"
	Max               StatName = "max"
	Mean              StatName = "mean"
	StandardDeviation StatName = "sd"
	Variance          StatName = "variance"

	DefaultPrecision = 3
	NumberBase       = 10
)

type Stats[T Number] interface {
	Push(T)
	Clear()
	Mean(format bool) float64
	Variance(format bool) float64
	StandardDeviation(format bool) float64
	Min(format bool) T
	Max(format bool) T
	Count(format bool) uint64
	AllStats(format bool) map[StatName]T
	ToJSON() string
	FromJSON(jsonStr string) error
}

type Stat[T Number] struct {
	count       uint64
	oldMean     float64
	newMean     float64
	oldVariance float64
	newVariance float64
	maximum     T
	minimum     T
	config      Config
	fmtModifier float64
}

type Config struct {
	Precision uint `json:"precision"`
}

func NewStat[T Number]() Stats[T] {
	s := Stat[T]{}

	s.SetPrecision(DefaultPrecision)

	return &s
}

func (s *Stat[T]) SetPrecision(precision uint) {
	s.config.Precision = precision
	s.fmtModifier = math.Pow(NumberBase, float64(precision))
}

func (s *Stat[T]) load(state map[string]interface{}) {
	for k, v := range state {
		switch k {
		case "count":
			s.count = uint64(v.(T))
		case "oldMean":
			s.oldMean = v.(float64)
		case "newMean":
			s.newMean = v.(float64)
		case "oldVariance":
			s.oldVariance = v.(float64)
		case "newVariance":
			s.newVariance = v.(float64)
		case "maximum":
			switch any(s.maximum).(type) {
			case uint64:
				m := v.(uint64)
				s.maximum = T(m)
			case float64:
				m := v.(float64)
				s.maximum = T(m)
			}
		case "minimum":
			switch any(s.minimum).(type) {
			case uint64:
				m := v.(uint64)
				s.minimum = T(m)
			case float64:
				m := v.(float64)
				s.minimum = T(m)
			}
		case "config":
			for k, cv := range v.(map[string]any) {
				switch k {
				case "precision":
					s.config.Precision = uint(cv.(T))
				}
			}
		}
	}

	if s.config.Precision == 0 {
		s.config.Precision = DefaultPrecision
	}

	s.SetPrecision(s.config.Precision)
}

func (s *Stat[T]) Push(num T) {
	s.count++

	if s.count == 1 {
		s.oldMean = float64(num)
		s.newMean = float64(num)
		s.oldVariance = 0.0
		s.minimum = num
		s.maximum = num

		return
	}

	s.newMean = s.oldMean + (float64(num)-s.oldMean)/float64(s.count)
	s.newVariance = s.oldVariance + (float64(num)-s.oldMean)*(float64(num)-s.newMean)
	s.minimum = T(math.Min(float64(s.minimum), float64(num)))
	s.maximum = T(math.Max(float64(s.maximum), float64(num)))

	// setup for next push
	s.oldMean = s.newMean
	s.oldVariance = s.newVariance
}

func (s *Stat[T]) Clear() {
	s.count = 0
	s.oldMean = 0.0
	s.newMean = 0.0
	s.oldVariance = 0.0
	s.newVariance = 0.0

	switch any(s.maximum).(type) {
	case uint64:
		s.maximum = 0
		s.minimum = 0
	case float64:
		s.maximum = 0.0
		s.minimum = 0.0
	}
}

func (s *Stat[T]) Mean(format bool) float64 {
	return float64(s.format(T(s.newMean), format))
}

func (s *Stat[T]) Min(format bool) T {
	return s.format(s.minimum, format)
}

func (s *Stat[T]) Max(format bool) T {
	return s.format(s.maximum, format)
}

func (s *Stat[T]) Count(format bool) uint64 {
	return uint64(s.format(T(s.count), format))
}

func (s *Stat[T]) StandardDeviation(format bool) float64 {
	return float64(s.format(T(math.Sqrt(s.Variance(false))), format))
}

func (s *Stat[T]) Variance(format bool) float64 {
	if s.count < 2 {
		return 0.0
	}

	return float64(s.format(T(s.newVariance/float64(s.count-1)), format))
}

func (s *Stat[T]) AllStats(format bool) map[StatName]T {
	return map[StatName]T{
		Count:             T(s.Count(format)),
		Min:               s.Min(format),
		Max:               s.Max(format),
		Mean:              T(s.Mean(format)),
		StandardDeviation: T(s.StandardDeviation(format)),
		Variance:          T(s.Variance(format)),
	}
}

func (s *Stat[T]) ToJSON() string {
	b, err := s.MarshalJSON()

	if err != nil {
		return ""
	}

	return string(b)
}

func (s *Stat[T]) FromJSON(jsonStr string) error {
	return s.UnmarshalJSON([]byte(jsonStr))
}

func (s *Stat[T]) format(val T, format bool) T {
	if !format {
		// don't format the value, just return it
		return val
	}

	return T(math.Round(float64(val)*s.fmtModifier) / s.fmtModifier)
}

func (s *Stat[T]) MarshalJSON() ([]byte, error) {
	o := map[string]interface{}{
		"config":      s.config,
		"count":       stringToNumber[uint64](strconv.FormatUint(s.count, NumberBase)),
		"oldMean":     stringToNumber[float64](strconv.FormatFloat(s.oldMean, 'f', -1, 64)),
		"newMean":     stringToNumber[float64](strconv.FormatFloat(s.newMean, 'f', -1, 64)),
		"oldVariance": stringToNumber[float64](strconv.FormatFloat(s.oldVariance, 'f', -1, 64)),
		"newVariance": stringToNumber[float64](strconv.FormatFloat(s.newVariance, 'f', -1, 64)),
	}

	switch any(s.maximum).(type) {
	case uint64:
		o["maximum"] = stringToNumber[uint64](strconv.FormatUint(uint64(s.maximum), NumberBase))
		o["minimum"] = stringToNumber[uint64](strconv.FormatUint(uint64(s.minimum), NumberBase))
	case float64:
		o["maximum"] = stringToNumber[float64](strconv.FormatFloat(float64(s.maximum), 'f', -1, 64))
		o["minimum"] = stringToNumber[float64](strconv.FormatFloat(float64(s.minimum), 'f', -1, 64))
	}

	return json.Marshal(o)
}

func (s *Stat[T]) UnmarshalJSON(bytes []byte) error {
	o := make(map[string]any)

	if err := json.Unmarshal(bytes, &o); err != nil {
		return err
	}

	s.load(o)

	return nil
}

func stringToNumber[T Number](num string) T {
	var n T

	switch any(n).(type) {
	case uint64:
		u, err := strconv.ParseUint(num, NumberBase, 64)

		if err != nil {
			break
		}

		n = T(u)
	case float64:
		f, err := strconv.ParseFloat(num, 64)

		if err != nil {
			break
		}

		n = T(f)
	}

	return n
}
