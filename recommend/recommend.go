// Package recommend provides recommend for closing channels based on the config
// provided. Only open public channels that have been monitored for some sane
// period of time will be considered for closing.
//
// Channels will be assessed based on the following data points:
// - Uptime percentage
package recommend

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/terminator/dataset"
	"github.com/lightningnetwork/lnd/lnrpc"
)

// errZeroMinimumAge is returned when the minimum ages provided by the config
// is zero.
var errZeroMinimumAge = errors.New("must provide a non-zero minimum age for " +
	"channel exclusion")

// CloseRecommendationConfig provides the functions and parameters required to
// provide close recommendations.
type CloseRecommendationConfig struct {
	// OpenChannels is a function which returns all of our currently open,
	// public channels.
	OpenChannels func() ([]*lnrpc.Channel, error)

	// StrongOutlier determines how much of an outlier a data point must be to
	// be considered an outlier.
	StrongOutlier bool

	// MinimumAge is the minimum amount of time that a channel must have been
	// monitored for before it is considered for closing.
	MinimumAge time.Duration
}

// Report contains a set of close recommendations and information about the
// number of channels considered for close.
type Report struct {
	// TotalChannels is the number of channels that we have.
	TotalChannels int

	// ConsideredChannels is the number of channels that have been monitored
	// for long enough to be considered for close.
	ConsideredChannels int

	// Recommendations is a map of chanel outpoints to a bool which indicates
	// whether we should close the channel.
	Recommendations map[wire.OutPoint]bool
}

// CloseRecommendations returns a Report which contains information about the
// channels that were considered and a list of close recommendations. Channels
// are considered for close if their uptime percentage is a lower outlier in
// uptime percentage dataset.
func CloseRecommendations(cfg CloseRecommendationConfig) (*Report, error) {
	// Check that the minimum wait time is non-zero.
	if cfg.MinimumAge == 0 {
		return nil, errZeroMinimumAge
	}

	stats, err := getChannelStats(cfg.OpenChannels, cfg.MinimumAge)
	if err != nil {
		return nil, err
	}

	recs := getCloseRecs(
		stats.channelUptimes, func(uptime float64) dataset.OutlierResult {
			return stats.uptimeDataset.IsOutlier(uptime, cfg.StrongOutlier)
		},
	)

	return &Report{
		TotalChannels:      stats.total,
		ConsideredChannels: len(stats.channelUptimes),
		Recommendations:    recs,
	}, nil
}

// getCloseRecs generates map of channel outpoint to bools indicating whether we
// should close the channel. It takes a map of channel points to their uptime
// percentage and a function which classifies a value as an outlier.
func getCloseRecs(channelUptimes map[wire.OutPoint]float64,
	uptimeOutlier func(uptime float64) dataset.OutlierResult) map[wire.OutPoint]bool {

	recommendations := make(map[wire.OutPoint]bool)

	for channel, uptime := range channelUptimes {
		// Check whether the channel is an outlier within the uptime dataset.
		outlier := uptimeOutlier(uptime)

		// If the channel is a lower outlier, penalize it.
		if outlier.LowerOutlier {
			recommendations[channel] = true
		}

		// If the channel is an upper outlier, reward it.
		if outlier.UpperOutlier {
			recommendations[channel] = false
		}
	}

	return recommendations
}

// channelStats contains a map of channels to their individual scores and the
// datasets that the scores are a part of.
type channelStats struct {
	// total is the total number of channels that are open.
	total int

	// channelUptimes is a map of channels to uptime percentage. This map only
	// contains channels which have been monitored for the required minimum
	// time.
	channelUptimes map[wire.OutPoint]float64

	// uptimeDataset is a dataset which contains the uptime values for all of the
	// channels which are being considered for close. It will be used to determine
	// which channels are outliers within the set.
	uptimeDataset *dataset.Dataset
}

// getChannelStats queries for the set of currently open channels and produces
// a channelStats struct containing a map of individual channels to their
// uptime percentage and a dataset containing all uptime percentages. Note that
// channels will be filtered out if they have not been monitored for long
// enough.
func getChannelStats(openChannels func() ([]*lnrpc.Channel, error),
	waitTime time.Duration) (*channelStats, error) {

	// Get the set of open channels that we are evaluating.
	openChans, err := openChannels()
	if err != nil {
		return nil, err
	}

	var (
		// Create a map which will hold channel point to uptime percentage.
		channels = make(map[wire.OutPoint]float64)

		// Accumulate the uptime scores for each channel so we can create
		// a dataset for them.
		uptimeData = make([]float64, len(openChans))
	)

	for _, channel := range openChans {
		outpoint, err := getOutPointFromString(channel.ChannelPoint)
		if err != nil {
			return nil, err
		}

		if channel.Lifetime < int64(waitTime.Seconds()) {
			log.Tracef("Channel: %v has not been monitored for long "+
				"enough, excluding it from consideration", channel.ChannelPoint)
			continue
		}

		// Calculate the uptime percentage for the channel and add it to the
		// channel -> uptime map.
		uptimePercentage := float64(channel.Uptime) / float64(channel.Lifetime)
		channels[*outpoint] = uptimePercentage

		log.Tracef("channel: %v has uptime percentage: %v", outpoint,
			uptimePercentage)

		// Add the uptime percentage to our set of uptime values.
		uptimeData = append(uptimeData, uptimePercentage)
	}

	log.Debugf("considering: % channels for close out of %v",
		len(channels), len(openChans))

	// Create a dataset for the uptime values we have collected.
	uptimeDataset, err := dataset.New(uptimeData)
	if err != nil {
		return nil, err
	}

	// Return channel stats which contains a mapping of channels to their
	// uptime and the dataset which contains all uptime values.
	return &channelStats{
		total:          len(openChans),
		channelUptimes: channels,
		uptimeDataset:  uptimeDataset,
	}, nil
}

// getOutPointFromString gets the channel outpoint from a string.
func getOutPointFromString(chanStr string) (*wire.OutPoint, error) {
	chanpoint := strings.Split(chanStr, ":")
	if len(chanpoint) != 2 {
		return nil, fmt.Errorf("expected 2 parts of channel point, "+
			"got: %v", len(chanpoint))
	}

	index, err := strconv.ParseInt(chanpoint[1], 10, 32)
	if err != nil {
		return nil, err
	}

	hash, err := chainhash.NewHashFromStr(chanpoint[0])
	if err != nil {
		return nil, err
	}

	return &wire.OutPoint{
		Hash:  *hash,
		Index: uint32(index),
	}, nil
}
