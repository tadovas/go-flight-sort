package flight_sorter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoFlightsErrorReturnedIfEmptyInput(t *testing.T) {
	_, err := SortFlights([]Flight{})
	assert.Equal(t, ErrNoFlights, err)
}

func TestConnectedFlightIsReturnedForMultipleFlights(t *testing.T) {
	tcs := []struct {
		input    []Flight
		expected Flight
	}{
		{
			input: []Flight{
				{Source: "SFO", Dest: "EWR"},
			},
			expected: Flight{Source: "SFO", Dest: "EWR"},
		},
		{
			input: []Flight{
				{Source: "ATL", Dest: "EWR"},
				{Source: "SFO", Dest: "ATL"},
			},
			expected: Flight{Source: "SFO", Dest: "EWR"},
		},
		{
			input: []Flight{
				{Source: "IND", Dest: "EWR"},
				{Source: "SFO", Dest: "ATL"},
				{Source: "GSO", Dest: "IND"},
				{Source: "ATL", Dest: "GSO"},
			},
			expected: Flight{Source: "SFO", Dest: "EWR"},
		},
	}

	for i, tc := range tcs {
		actual, err := SortFlights(tc.input)
		assert.NoError(t, err, "test case %v", i)
		assert.Equal(t, tc.expected, actual, "test case %v unexpected result", i)
	}

}

func TestFlightLoopIsDetected(t *testing.T) {
	_, err := SortFlights([]Flight{
		{Source: "A", Dest: "B"},
		{Source: "B", Dest: "A"},
	})
	assert.Equal(t, ErrFlightLoop, err)
}

func TestBrokenFlightPathIsDetected(t *testing.T) {
	_, err := SortFlights([]Flight{
		{Source: "A", Dest: "B"},
		// missing: "B" -> "C"
		{Source: "C", Dest: "D"},
	})
	assert.Equal(t, ErrBrokenPath, err)
}

func TestSplitPathIsDetected(t *testing.T) {
	_, err := SortFlights([]Flight{
		{Source: "A", Dest: "B"},
		{Source: "A", Dest: "C"},
	})
	assert.Equal(t, ErrBrokenPath, err)
}

func TestMergedPathIsDetected(t *testing.T) {
	_, err := SortFlights([]Flight{
		{Source: "A", Dest: "B"},
		{Source: "C", Dest: "B"},
	})
	assert.Equal(t, ErrBrokenPath, err)
}

const testJson = `
{
	"source" : "A",
	"dest" : "B"
}
`

func TestFlightIsUnmarshaledFromJson(t *testing.T) {
	var flight Flight
	assert.NoError(t, json.Unmarshal([]byte(testJson), &flight))
	assert.Equal(t, Flight{
		Source: "A",
		Dest:   "B",
	}, flight)
}
