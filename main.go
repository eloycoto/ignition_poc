package main

import (

	// "github.com/eloycoto/ignition_poc/pkg/ignition_config"
	"log"

	"github.com/eloycoto/ignition_poc/pkg/ignition"
	_ "github.com/eloycoto/ignition_poc/pkg/ignition/source/exec/stages/files"
)

var (
	sampleConfig = `
{
  "ignition": { "version": "3.0.0" },
  "storage": {
    "files": [{
      "path": "/opt/someconfig",
      "mode": 420,
      "contents": { "source": "data:,example%20file%0A" },
      "overwrite": true
    }]
  },
  "systemd": {
    "units": [{
      "name": "eloy.service",
      "enabled": true,
      "contents": "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"
    }]
  }
}
`
)

func main() {
	cfg, err := ignition.ParseConfig(sampleConfig)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	err = ignition.RunConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to run config: %v", err)
	}

	// fmt.Println("Hello world")
	// fmt.Printf("Sample config:\n %s \n", sampleConfig)

	// logger := log.New(true)
	// defer logger.Close()

	// state, err := state.Load("/tmp/ignition")
	// if err != nil {
	// 	logger.Crit("reading state: %s", err)
	// 	os.Exit(3)
	// }

	// platformConfig := platform.MustGet("file")
	// fetcher, err := platformConfig.NewFetcherFunc()(&logger)
	// if err != nil {
	// 	logger.Crit("failed to generate fetcher: %s", err)

	// 	os.Exit(3)
	// }

	// cfg, rpt, err := config.Parse([]byte(sampleConfig))
	// // logger.LogReport(rpt)
	// if rpt.IsFatal() || err != nil {
	// 	logger.Crit("couldn't parse config: %v", err)
	// 	os.Exit(1)
	// }

	// cfgFetcher := exec.ConfigFetcher{
	// 	Logger:  &logger,
	// 	Fetcher: &fetcher,
	// 	State:   &state,
	// }

	// finalCfg, err := cfgFetcher.RenderConfig(cfg)
	// if err != nil {
	// 	logger.Crit("Failed here: %v", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("Stages:: ", stages.Get("files"))
	// // fileStage := stages.Get("files")

	// stage := stages.Get("files").Create(&logger, "/", fetcher, &state)
	// if err := stage.Apply(finalCfg, true); err != nil {
	// 	logger.Crit("running stage 'files': %w", err)
	// }
	// // if err := fileStage.Apply(finalCfg, flags.IgnoreUnsupported); err != nil {
	// // 	return fmt.Errorf("running stage '%s': %w", stageName, err)
	// // }

	// // // logger := log.New(os.Stdout)
	// // // // defer logger.Close()

	// // if err := apply.Run(cfg, apply.Flags{}, nil); err != nil {
	// // 	// logger.Crit("failed to apply: %v", err)
	// // 	fmt.Println("Failed", err)
	// // 	os.Exit(1)
	// // }
	// fmt.Println("Phase1")
}
