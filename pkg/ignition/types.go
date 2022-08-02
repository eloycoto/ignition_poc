package ignition

import (
	"github.com/coreos/ignition/v2/config/v3_2/types"
)

type Config struct {
	Storage Storage       `json:"storage,omitempty"`
	Systemd types.Systemd `json:"systemd,omitempty"`
}

type Storage struct {
	Directories []types.Directory `json:"directories,omitempty"`
	Files       []types.File      `json:"files,omitempty"`
	Links       []types.Link      `json:"links,omitempty"`
}
