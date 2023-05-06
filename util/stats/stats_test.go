package stats

import "testing"

var testData = []map[string]any{
	{
		"data": []float64{},
		"results": map[string]float64{
			"mean":              0,
			"variance":          0,
			"standardDeviation": 0,
			"min":               0,
			"max":               0,
		},
	},
	{
		"data": []float64{1, 2, 3, 4, 5},
		"results": map[string]float64{
			"mean":              3,
			"variance":          2.5,
			"standardDeviation": 1.581,
			"min":               1,
			"max":               5,
		},
	},
}

func TestCalcStats(t *testing.T) {
	for _, test := range testData {
		stats := NewStat[float64]()

		data := test["data"].([]float64)

		for _, v := range data {
			stats.Push(v)
		}

		wanted := test["results"].(map[string]float64)

		if stats.Mean(false) != wanted["mean"] {
			t.Errorf("stats.Mean(false) = %v, wanted: %v", stats.Mean(false), wanted["mean"])
		}

		if stats.Variance(false) != wanted["variance"] {
			t.Errorf("stats.Variance(false) = %v, wanted: %v", stats.Variance(false), wanted["variance"])
		}

		if stats.StandardDeviation(true) != wanted["standardDeviation"] {
			t.Errorf("stats.StandardDeviation(true) = %v, wanted: %v", stats.StandardDeviation(true), wanted["standardDeviation"])
		}

		if stats.Min(false) != wanted["min"] {
			t.Errorf("stats.Min(false) = %v, wanted: %v", stats.Min(false), wanted["min"])
		}

		if stats.Max(false) != wanted["max"] {
			t.Errorf("stats.Max(false) = %v, wanted: %v", stats.Max(false), wanted["max"])
		}
	}
}

func TestStat_MarshalJSON_UnmarshalJSON(t *testing.T) {
	test := testData[1]
	expectedJson := `{"config":{"precision":3},"count":5,"maximum":5,"minimum":1,"newMean":3,"newVariance":10,"oldMean":3,"oldVariance":10}`

	stats := NewStat[float64]()

	data := test["data"].([]float64)

	for _, v := range data {
		stats.Push(v)
	}

	wanted := test["results"].(map[string]float64)

	if stats.Mean(false) != wanted["mean"] {
		t.Errorf("stats.Mean(false) = %v, wanted: %v", stats.Mean(false), wanted["mean"])
	}

	if stats.Variance(false) != wanted["variance"] {
		t.Errorf("stats.Variance(false) = %v, wanted: %v", stats.Variance(false), wanted["variance"])
	}

	if stats.StandardDeviation(true) != wanted["standardDeviation"] {
		t.Errorf("stats.StandardDeviation(true) = %v, wanted: %v", stats.StandardDeviation(true), wanted["standardDeviation"])
	}

	if stats.Min(false) != wanted["min"] {
		t.Errorf("stats.Min(false) = %v, wanted: %v", stats.Min(false), wanted["min"])
	}

	if stats.Max(false) != wanted["max"] {
		t.Errorf("stats.Max(false) = %v, wanted: %v", stats.Max(false), wanted["max"])
	}

	json := stats.ToJSON()

	if json != expectedJson {
		t.Errorf("stats.ToJSON() = '%v', wanted: '%v'", json, expectedJson)
	}

	stats = NewStat[float64]()

	if stats.Mean(false) == wanted["mean"] {
		t.Errorf("stats.Mean(false) = %v, wanted: != %v", stats.Mean(false), wanted["mean"])
	}

	if stats.Variance(false) == wanted["variance"] {
		t.Errorf("stats.Variance(false) = %v, wanted: != %v", stats.Variance(false), wanted["variance"])
	}

	if stats.StandardDeviation(true) == wanted["standardDeviation"] {
		t.Errorf("stats.StandardDeviation(true) = %v, wanted: != %v", stats.StandardDeviation(true), wanted["standardDeviation"])
	}

	if stats.Min(false) == wanted["min"] {
		t.Errorf("stats.Min(false) = %v, wanted: != %v", stats.Min(false), wanted["min"])
	}

	if stats.Max(false) == wanted["max"] {
		t.Errorf("stats.Max(false) = %v, wanted: != %v", stats.Max(false), wanted["max"])
	}

	if err := stats.FromJSON(json); err != nil {
		t.Errorf("stats.FromJSON() returned an error, error: %v", err)
	}

	if stats.Mean(false) != wanted["mean"] {
		t.Errorf("stats.Mean(false) = %v, wanted: %v", stats.Mean(false), wanted["mean"])
	}

	if stats.Variance(false) != wanted["variance"] {
		t.Errorf("stats.Variance(false) = %v, wanted: %v", stats.Variance(false), wanted["variance"])
	}

	if stats.StandardDeviation(true) != wanted["standardDeviation"] {
		t.Errorf("stats.StandardDeviation(true) = %v, wanted: %v", stats.StandardDeviation(true), wanted["standardDeviation"])
	}

	if stats.Min(false) != wanted["min"] {
		t.Errorf("stats.Min(false) = %v, wanted: %v", stats.Min(false), wanted["min"])
	}

	if stats.Max(false) != wanted["max"] {
		t.Errorf("stats.Max(false) = %v, wanted: %v", stats.Max(false), wanted["max"])
	}
}
