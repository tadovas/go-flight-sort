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
				{Start: "SFO", Dest: "EWR"},
			},
			expected: Flight{Start: "SFO", Dest: "EWR"},
		},
		{
			input: []Flight{
				{Start: "ATL", Dest: "EWR"},
				{Start: "SFO", Dest: "ATL"},
			},
			expected: Flight{Start: "SFO", Dest: "EWR"},
		},
		{
			input: []Flight{
				{Start: "IND", Dest: "EWR"},
				{Start: "SFO", Dest: "ATL"},
				{Start: "GSO", Dest: "IND"},
				{Start: "ATL", Dest: "GSO"},
			},
			expected: Flight{Start: "SFO", Dest: "EWR"},
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
		{Start: "A", Dest: "B"},
		{Start: "B", Dest: "A"},
	})
	assert.Equal(t, ErrFlightLoop, err)
}

func TestBrokenFlightPathIsDetected(t *testing.T) {
	_, err := SortFlights([]Flight{
		{Start: "A", Dest: "B"},
		// missing: "B" -> "C"
		{Start: "C", Dest: "D"},
	})
	assert.Equal(t, ErrBrokenPath, err)
}

func TestSplitPathIsDetected(t *testing.T) {
	_, err := SortFlights([]Flight{
		{Start: "A", Dest: "B"},
		{Start: "A", Dest: "C"},
	})
	assert.Equal(t, ErrBrokenPath, err)
}

func TestMergedPathIsDetected(t *testing.T) {
	_, err := SortFlights([]Flight{
		{Start: "A", Dest: "B"},
		{Start: "C", Dest: "B"},
	})
	assert.Equal(t, ErrBrokenPath, err)
}

const testJson = `
{
	"start" : "A",
	"dest" : "B"
}
`

func TestFlightIsUnmarshaledFromJson(t *testing.T) {
	var flight Flight
	assert.NoError(t, json.Unmarshal([]byte(testJson), &flight))
	assert.Equal(t, Flight{
		Start: "A",
		Dest:  "B",
	}, flight)
}
