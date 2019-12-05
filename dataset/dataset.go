// Package dataset provides a basic dataset type which calculates the
// median, upper and lower quartiles for float64 typed data.
package dataset

import (
	"fmt"
	"sort"
)

const (
	// weakOutlier is the multiplier that we apply to the inter-quartile range
	// to calculate weak outliers.
	weakOutlier = 1.5

	// strongOutlier is the multiplier that we apply to the inter-quartile range
	// to calculate strong outliers.
	strongOutlier = 3
)

var (
	// errNoValues is returned when an attempt is made to calculate the median of
	// a zero length array.
	errNoValues = fmt.Errorf("can't calculate median for zero length " +
		"array")

	// ErrTooFewValues is returned when there are too few values provided to
	// calculate quartiles.
	ErrTooFewValues = fmt.Errorf("can't calculate quartiles for fewer than 3 " +
		"elements")
)

// Dataset contains information about a set of float64 data points.
type Dataset struct {
	values        []float64
	UpperQuartile float64
	Median        float64
	LowerQuartile float64
}

// getMedian gets the median for a set of *already sorted* values. It returns
// an error if there are no values.
func getMedian(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, errNoValues
	}

	if len(values)%2 == 0 {
		split := len(values) / 2
		midpoint := (values[split-1] + values[split]) / 2
		return midpoint, nil
	}

	// If there is an odd number of values in the dataset, return the middle
	// element as the median, and the remaining values on either side.
	split := len(values) / 2
	return values[split], nil
}

// New returns a dataset with the median and upper/lower quartiles calculated.
// If there are fewer than three values in the dataset, an error is returned
// because we cannot calculate quartiles.
func New(values []float64) (*Dataset, error) {
	if len(values) < 3 {
		return nil, ErrTooFewValues
	}

	// Sort the dataset in ascending order.
	sort.Float64s(values)

	median, err := getMedian(values)
	if err != nil {
		return nil, err
	}

	// Get the cutoff points for calculating the lower and upper quartiles.
	// The "exclusive" method of calculating quartiles is used, meaning that
	// the dataset is split in half, excluding the median value in the case
	// of an odd number of elements.
	var cutoffLower, cutoffUpper int
	if len(values)%2 == 0 {
		// For an even number of elements, we split the dataset exactly in half.
		cutoffLower = len(values) / 2
		cutoffUpper = len(values) / 2
	} else {
		// For an odd number of elements, we exclude the middle element cut
		// the dataset in half on either side.
		cutoffLower = (len(values) - 1) / 2
		cutoffUpper = cutoffLower + 1
	}

	lowerQuartile, err := getMedian(values[:cutoffLower])
	if err != nil {
		return nil, err
	}

	upperQuartile, err := getMedian(values[cutoffUpper:])
	if err != nil {
		return nil, err
	}

	return &Dataset{
		values:        values,
		Median:        median,
		LowerQuartile: lowerQuartile,
		UpperQuartile: upperQuartile,
	}, nil
}

// interQuartileRange returns the inter-quartile range for the dataset.
func (d *Dataset) interQuartileRange() float64 {
	return d.UpperQuartile - d.LowerQuartile
}

// OutlierResult returns the results of an outlier check.
type OutlierResult struct {
	UpperOutlier bool
	LowerOutlier bool
}

// IsOutlier returns an outlier result which indicates whether a value is an
// upper or lower outlier (or not an outlier) for the dataset given. The strong
// bool is set to adjust whether we check for a strong weak outlier.
func (d *Dataset) IsOutlier(value float64, strong bool) OutlierResult {
	outlierMagintude := weakOutlier
	if strong {
		outlierMagintude = strongOutlier
	}

	lowerLimit := d.LowerQuartile - (d.interQuartileRange() * outlierMagintude)
	upperLimit := d.UpperQuartile + (d.interQuartileRange() * outlierMagintude)
	return OutlierResult{
		UpperOutlier: value > upperLimit,
		LowerOutlier: value < lowerLimit,
	}
}
