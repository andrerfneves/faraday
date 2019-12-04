package terminator

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightningnetwork/lnd/channeldb"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnwire"
)

func getRPCRevenue(ctx context.Context, client lnrpc.LightningClient) (*RevenueReport, error) {
	channels, err := client.ListChannels(ctx, &lnrpc.ListChannelsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error calling list channels: %v", err)
	}

	closed, err := client.ClosedChannels(ctx, &lnrpc.ClosedChannelsRequest{})
	if err != nil {
		return nil, err
	}

	// Add the channels looked up to a map of short channel id to outpoint.
	channelIDs := make(map[lnwire.ShortChannelID]wire.OutPoint)
	for _, channel := range channels.Channels {
		outPoint, err := getChanPoint(channel.ChannelPoint)
		if err != nil {
			return nil, err
		}

		channelIDs[lnwire.NewShortChanIDFromInt(channel.ChanId)] = *outPoint
	}

	for _, closedChannel := range closed.Channels {
		outPoint, err := getChanPoint(closedChannel.ChannelPoint)
		if err != nil {
			return nil, err
		}

		channelIDs[lnwire.NewShortChanIDFromInt(closedChannel.ChanId)] = *outPoint
	}

	query := func(offset, maxEvents uint32) (channeldb.ForwardingLogTimeSlice, error) {
		start := time.Now().Add(time.Hour * -24)
		end := time.Now()

		req := &lnrpc.ForwardingHistoryRequest{
			StartTime:    uint64(start.Unix()),
			EndTime:      uint64(end.Unix()),
			IndexOffset:  offset,
			NumMaxEvents: maxEvents,
		}

		resp, err := client.ForwardingHistory(ctx, req)
		if err != nil {
			return channeldb.ForwardingLogTimeSlice{}, err
		}

		var events []channeldb.ForwardingEvent
		for _, e := range resp.ForwardingEvents {
			events = append(events, channeldb.ForwardingEvent{
				Timestamp:      time.Unix(int64(e.Timestamp), 0),
				IncomingChanID: lnwire.NewShortChanIDFromInt(e.ChanIdIn),
				OutgoingChanID: lnwire.NewShortChanIDFromInt(e.ChanIdOut),
				AmtIn:          lnwire.MilliSatoshi(e.AmtIn),
				AmtOut:         lnwire.MilliSatoshi(e.AmtOut),
			})
		}
		return channeldb.ForwardingLogTimeSlice{
			ForwardingEventQuery: channeldb.ForwardingEventQuery{
				StartTime:    start,
				EndTime:      end,
				IndexOffset:  req.IndexOffset,
				NumMaxEvents: req.NumMaxEvents,
			},
			ForwardingEvents: events,
			LastIndexOffset:  resp.LastOffsetIndex,
		}, nil
	}

	// Query for report over relevant period.
	return GetRevenueReport(channelIDs, query, 0.5)
}

// getChanPoint gets the channel outpoint from a string.
func getChanPoint(chanStr string) (*wire.OutPoint, error) {
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
