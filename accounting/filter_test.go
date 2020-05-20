package accounting

import (
	"testing"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/stretchr/testify/require"
)

var (
	startTime   int64 = 100000
	inRangeTime int64 = 200000
	endTime     int64 = 300000
)

// TestInRange tests filtering of timestamps by a inclusive start time and
// exclusive end time.
func TestInRange(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		inRange   bool
	}{
		{
			name:      "before start time - not in range",
			timestamp: startTime - 100,
			inRange:   false,
		},
		{
			name:      "equals start time - ok",
			timestamp: startTime,
			inRange:   true,
		},
		{
			name:      "between start and end - ok",
			timestamp: inRangeTime,
			inRange:   true,
		},
		{
			name:      "equals end time - not in range",
			timestamp: endTime,
			inRange:   false,
		},
		{
			name:      "after end time - not in range",
			timestamp: endTime + 10,
			inRange:   false,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			start := time.Unix(startTime, 0)
			end := time.Unix(endTime, 0)
			ts := time.Unix(test.timestamp, 0)

			inRange := inRange(ts, start, end)
			require.Equal(t, test.inRange, inRange)
		})
	}
}

// TestFilterOnChain tests filtering transactions based on timestamp and
// confirmations.
func TestFilterOnChain(t *testing.T) {
	// Create three test transactions, one confirmed but outside of our
	// range, one confirmed and in our range and one in our range but not
	// confirmed.
	confirmedTxOutOfRange := &lnrpc.Transaction{
		TimeStamp:        startTime - 10,
		NumConfirmations: 1,
	}

	confirmedTx := &lnrpc.Transaction{
		TimeStamp:        inRangeTime,
		NumConfirmations: 1,
	}

	noConfTx := &lnrpc.Transaction{
		TimeStamp:        inRangeTime,
		NumConfirmations: 0,
	}

	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	unfiltered := []*lnrpc.Transaction{
		confirmedTx, noConfTx, confirmedTxOutOfRange,
	}
	filtered := filterOnChain(start, end, unfiltered)

	// We only expect our confirmed transaction in the time range we
	// specified to be included.
	expected := []*lnrpc.Transaction{confirmedTx}
	require.Equal(t, expected, filtered)
}