package dataset

import "testing"

// TestGetMedian tests median calculation for a series of inputs, including
// the error case where there are no values.
func TestGetMedian(t *testing.T) {
	tests := []struct {
		name           string
		values         []float64
		expectedErr    error
		expectedMedian float64
	}{
		{
			name:        "no values",
			values:      []float64{},
			expectedErr: errNoValues,
		},
		{
			name:           "one value",
			values:         []float64{1},
			expectedErr:    nil,
			expectedMedian: 1,
		},
		{
			name:           "two values",
			values:         []float64{1, 2},
			expectedErr:    nil,
			expectedMedian: 1.5,
		},
		{
			name:           "three values",
			values:         []float64{1, 2, 3},
			expectedErr:    nil,
			expectedMedian: 2,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			median, err := getMedian(test.values)
			if err != test.expectedErr {
				t.Fatalf("expected: %v, got: %v", test.expectedErr, err)
			}

			if test.expectedMedian != median {
				t.Fatalf("expected: %v, got: %v", test.expectedMedian, median)
			}
		})
	}
}

// TestNew tests the creation of a new dataset from a set of values. It tests
// the case where the dataset does not have enough values, and cases with odd
// and even numbers of values to test the splitting of the dataset.
func TestNew(t *testing.T) {
	tests := []struct {
		name                  string
		values                []float64
		expectedErr           error
		expectedMedian        float64
		expectedLowerQuartile float64
		expectedUpperQuartile float64
	}{
		{
			name:        "no elements",
			values:      []float64{},
			expectedErr: ErrTooFewValues,
		},
		{
			name:                  "three elements",
			values:                []float64{3, 1, 2},
			expectedLowerQuartile: 1,
			expectedMedian:        2,
			expectedUpperQuartile: 3,
		},
		{
			name:                  "four elements",
			values:                []float64{1, 2, 3, 4},
			expectedLowerQuartile: 1.5,
			expectedMedian:        2.5,
			expectedUpperQuartile: 3.5,
		},
		{
			name:                  "five elements",
			values:                []float64{1, 2, 4, 3, 5},
			expectedLowerQuartile: 1.5,
			expectedMedian:        3,
			expectedUpperQuartile: 4.5,
		},
		{
			name:                  "four elements",
			values:                []float64{1, 2, 3, 4, 5, 6, 7, 8},
			expectedLowerQuartile: 2.5,
			expectedMedian:        4.5,
			expectedUpperQuartile: 6.5,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			dataset, err := New(test.values)
			if err != test.expectedErr {
				t.Fatalf("expected: %v, got: %v", test.expectedErr, err)
			}

			// If an error occurred, we do not need to perform any further checks.
			if err != nil {
				return
			}

			if test.expectedMedian != dataset.Median {
				t.Fatalf("expected: %v, got: %v",
					test.expectedMedian, dataset.Median)
			}

			if test.expectedLowerQuartile != dataset.LowerQuartile {
				t.Fatalf("expected: %v, got: %v",
					test.expectedLowerQuartile, dataset.LowerQuartile)
			}

			if test.expectedUpperQuartile != dataset.UpperQuartile {
				t.Fatalf("expected: %v, got: %v",
					test.expectedUpperQuartile, dataset.UpperQuartile)
			}
		})
	}
}
