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
}
