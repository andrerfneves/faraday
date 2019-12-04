// Package terminator contains the main function for the terminator.
package terminator

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/lightninglabs/loop/lndclient"
)

// Main is the real entry point for terminator. It is required to ensure that
// defers are properly executed when os.Exit() is called.
func Main() error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// NewBasicClient get a lightning rpc client with
	client, err := lndclient.NewBasicClient(
		config.RPCServer,
		config.TLSCertPath,
		config.MacaroonDir,
		config.network,
		lndclient.MacFilename(config.MacaroonFile),
	)
	if err != nil {
		return fmt.Errorf("cannot connect to lightning client: %v", err)
	}

	// Query for report over relevant period.
	report, err := getRPCRevenue(ctx, client)
	if err != nil {
		return err
	}

	spew.Dump(report)

	return nil
}
