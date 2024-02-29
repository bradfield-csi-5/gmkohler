package iter

import (
	"pg/tuple"
	"testing"
)

func TestTuple_GetColumnValue(t *testing.T) {
	var testData = []struct {
		colName       string
		expectedValue tuple.ColumnValue
		errorExpected bool
	}{
		{
			colName:       "name",
			expectedValue: tuple.ColumnValue("Johnny Bench"),
		},
		{
			colName:       "team",
			errorExpected: true,
		},
	}

	var tup = tuple.Tuple{
		Columns: []tuple.Column{
			{"name", "Johnny Bench"},
			{"position", "C"},
		},
	}

	for _, test := range testData {
		val, err := tup.GetColumnValue(test.colName)
		if err == nil {
			if test.errorExpected {
				t.Errorf("Expected failure, got %v\n", val)
			}
		} else if val != test.expectedValue {
			t.Errorf("Expected %v, got %v", test.expectedValue, val)
		}
	}
}
