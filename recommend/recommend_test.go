package recommend

import (
	"errors"
	"testing"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/terminator/dataset"
	"github.com/lightningnetwork/lnd/lnrpc"
)

// TestCloseRecommendations tests CloseRecommendations for error cases where
// the function provided to list channels fails or the config provided is
// invalid. It also has cases for calls which return not enough channels, and
// the minimum acceptable number of channels. It does not test the report
// provided, because that will be covered by further tests.
func TestCloseRecommendations(t *testing.T) {
	var openChanErr = errors.New("intentional test err")

	tests := []struct {
		name         string
		OpenChannels func() ([]*lnrpc.Channel, error)
		MinAge       time.Duration
		expectedErr  error
	}{
		{
			name:         "no channels",
			OpenChannels: func() ([]*lnrpc.Channel, error) { return nil, nil },
			MinAge:       time.Hour,
			expectedErr:  dataset.ErrTooFewValues,
		},
		{
			name:         "open channels fails",
			OpenChannels: func() ([]*lnrpc.Channel, error) { return nil, openChanErr },
			MinAge:       time.Hour,
			expectedErr:  openChanErr,
		},
		{
			name:         "zero min age",
			OpenChannels: func() ([]*lnrpc.Channel, error) { return nil, nil },
			MinAge:       0,
			expectedErr:  errZeroMinimumAge,
		},
		{
			name: "enough channels",
			OpenChannels: func() ([]*lnrpc.Channel, error) {
				return []*lnrpc.Channel{
					{ChannelPoint: "a:1"},
					{ChannelPoint: "b:2"},
					{ChannelPoint: "c:3"},
				}, nil
			},
			MinAge:      time.Hour,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := CloseRecommendations(CloseRecommendationConfig{
				OpenChannels:  test.OpenChannels,
				StrongOutlier: true,
				MinimumAge:    test.MinAge,
			})
			if err != test.expectedErr {
				t.Fatalf("expected: %v, got: %v", test.expectedErr, err)
			}
		})
	}
}

// TestGetCloseRecs tests the generating of close recommendations for a set of
// channels.
func TestGetCloseRecs(t *testing.T) {
	var (
		chan1 = wire.OutPoint{Index: 1}
		chan2 = wire.OutPoint{Index: 2}
	)

	tests := []struct {
		name           string
		channelUptimes map[wire.OutPoint]float64
		uptimeFunc     func(value float64) dataset.OutlierResult
		expectedRecs   map[wire.OutPoint]bool
	}{
		{
			name: "no outliers",
			channelUptimes: map[wire.OutPoint]float64{
				chan1: 0.5,
				chan2: 0.5,
			},
			uptimeFunc: func(value float64) dataset.OutlierResult {
				return dataset.OutlierResult{}
			},
			expectedRecs: map[wire.OutPoint]bool{
				chan1: false,
				chan2: false,
			},
		},
		{
			name: "upper outlier not included",
			channelUptimes: map[wire.OutPoint]float64{
				chan1: 0.5,
				chan2: 0.5,
			},
			uptimeFunc: func(value float64) dataset.OutlierResult {
				return dataset.OutlierResult{UpperOutlier: true}
			},
			expectedRecs: map[wire.OutPoint]bool{
				chan1: false,
				chan2: false,
			},
		},
		{
			name: "lower outlier not included",
			channelUptimes: map[wire.OutPoint]float64{
				chan1: 0.1,
				chan2: 0.5,
			},
			uptimeFunc: func(value float64) dataset.OutlierResult {
				if value == 0.1 {
					return dataset.OutlierResult{LowerOutlier: true}
				}
				return dataset.OutlierResult{}
			},
			expectedRecs: map[wire.OutPoint]bool{
				chan1: true,
				chan2: false,
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			recs := getCloseRecs(test.channelUptimes, test.uptimeFunc)

			// Run through our expected set of recommendations and check that they
			// match the set returned in the report.
			for channel, expectClose := range test.expectedRecs {
				recClose := recs[channel]
				if recClose != expectClose {
					t.Fatalf("expected close rec: %v for channel: %v, got: %v",
						expectClose, channel, recClose)
				}
			}
		})
	}
}

// TestGetChannelStats tests inclusion of channels based on their lifetime.
func TestGetChannelStats(t *testing.T) {
	openChans := func() (channels []*lnrpc.Channel, e error) {
		return []*lnrpc.Channel{
			{
				ChannelPoint: "a:1",
				Lifetime:     10,
				Uptime:       1,
			},
			{
				ChannelPoint: "a:2",
				Lifetime:     100,
				Uptime:       1,
			},
			{
				ChannelPoint: "a:3",
				Lifetime:     1000,
				Uptime:       1,
			},
		}, nil
	}

	tests := []struct {
		name             string
		openChannels     func() ([]*lnrpc.Channel, error)
		minAge           time.Duration
		expectedChannels int
	}{
		{
			name:             "channel not monitored for long enough",
			openChannels:     openChans,
			minAge:           time.Second * 10000,
			expectedChannels: 0,
		},
		{
			name:             "channel monitored for long enough",
			openChannels:     openChans,
			minAge:           time.Second * 100,
			expectedChannels: 1,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stats, err := getChannelStats(test.openChannels, test.minAge)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(stats.channelUptimes) != test.expectedChannels {
				t.Fatalf("expected: %v channels, got: %v",
					test.expectedChannels, len(stats.channelUptimes))
			}
		})
	}
}
