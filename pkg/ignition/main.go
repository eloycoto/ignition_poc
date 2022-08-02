package ignition

import (
	"fmt"

	// "github.com/eloycoto/ignition_poc/pkg/ignition_config"

	"github.com/eloycoto/ignition_poc/pkg/ignition/source/exec"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/exec/stages"
	_ "github.com/eloycoto/ignition_poc/pkg/ignition/source/exec/stages/files"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/log"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/platform"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/state"

	types_exp "github.com/coreos/ignition/v2/config/v3_4_experimental/types"

	"github.com/coreos/ignition/v2/config"
)

func ParseConfig(rawConfig string) (*types_exp.Config, error) {
	cfg, rpt, err := config.Parse([]byte(rawConfig))
	if rpt.IsFatal() || err != nil {
		return nil, err
	}
	return &cfg, nil
}

func RunConfig(cfg *types_exp.Config) error {

	logger := log.New(true)
	defer logger.Close()

	state, err := state.Load("/tmp/flotta_ignition")
	if err != nil {
		return fmt.Errorf("Cannot read state: %v", err)
	}

	platformConfig := platform.MustGet("file")
	fetcher, err := platformConfig.NewFetcherFunc()(&logger)
	if err != nil {
		return fmt.Errorf("Failed to generate fetcher: %v", err)
	}

	cfgFetcher := exec.ConfigFetcher{
		Logger:  &logger,
		Fetcher: &fetcher,
		State:   &state,
	}

	finalCfg, err := cfgFetcher.RenderConfig(*cfg)
	if err != nil {
		return fmt.Errorf("Failed on config reader: %v", err)
	}

	stage := stages.Get("files").Create(&logger, "/", fetcher, &state)
	if err := stage.Apply(finalCfg, true); err != nil {
		return fmt.Errorf("Running stage files failed: %v", err)
	}
	return nil
	// fmt.Println("Phase1")
}

