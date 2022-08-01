package main

import (
	"fmt"

	// "github.com/eloycoto/ignition_poc/pkg/ignition_config"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/exec/stages"
)

var (
	sampleConfig = `
{
  "ignition": { "version": "3.0.0" },
  "storage": {
    "files": [{
      "path": "/etc/someconfig",
      "mode": 420,
      "contents": { "source": "data:,example%20file%0A" }
    }]
  }
}
`
)

func main() {
	fmt.Println("Hello world")
	fmt.Printf("Sample config:\n %s \n", sampleConfig)
	// fmt.Println(ignition_config.Merda())

	// parsedData, err := ignition_config.Parse([]byte(sampleConfig))
	// fmt.Printf("ParsedData:\n %+v \n", parsedData)
	// fmt.Printf("Err:\n%s \n", err)

	fmt.Println("Stages:: ", stages.Get("files"))

}
