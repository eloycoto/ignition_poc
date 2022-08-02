// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The vmware provider fetches a configuration from the VMware Guest Info
// interface.

package vmware

import (
	"fmt"

	"github.com/coreos/ignition/v2/config/v3_4_experimental/types"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/providers"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/providers/util"
	"github.com/eloycoto/ignition_poc/pkg/ignition/source/resource"

	"github.com/coreos/vcontext/report"
	"github.com/vmware/vmw-guestinfo/rpcvmx"
	"github.com/vmware/vmw-guestinfo/vmcheck"
)

const (
	GUESTINFO_OVF               = "ovfenv"
	GUESTINFO_USERDATA          = "ignition.config.data"
	GUESTINFO_USERDATA_ENCODING = "ignition.config.data.encoding"

	GUESTINFO_DELETED_USERDATA          = "e30K"
	GUESTINFO_DELETED_USERDATA_ENCODING = "base64"

	OVF_PREFIX            = "guestinfo."
	OVF_USERDATA          = OVF_PREFIX + GUESTINFO_USERDATA
	OVF_USERDATA_ENCODING = OVF_PREFIX + GUESTINFO_USERDATA_ENCODING
)

func FetchConfig(f *resource.Fetcher) (types.Config, report.Report, error) {
	if isVM, err := vmcheck.IsVirtualWorld(true); err != nil {
		return types.Config{}, report.Report{}, err
	} else if !isVM {
		return types.Config{}, report.Report{}, providers.ErrNoProvider
	}

	config, err := fetchRawConfig(f)
	if err != nil {
		return types.Config{}, report.Report{}, err
	}

	decodedData, err := decodeConfig(config)
	if err != nil {
		f.Logger.Debug("failed to decode config: %v", err)
		return types.Config{}, report.Report{}, err
	}

	f.Logger.Debug("config successfully fetched")
	return util.ParseConfig(f.Logger, decodedData)
}

func fetchRawConfig(f *resource.Fetcher) (config, error) {
	info := rpcvmx.NewConfig()

	var ovfData string
	var ovfEncoding string

	ovfEnv, err := info.String(GUESTINFO_OVF, "")
	if err != nil {
		f.Logger.Warning("failed to fetch ovfenv: %v. Continuing...", err)
	} else if ovfEnv != "" {
		f.Logger.Debug("using OVF environment from guestinfo")
		env, err := ReadOvfEnvironment([]byte(ovfEnv))
		if err != nil {
			f.Logger.Warning("failed to parse OVF environment: %v. Continuing...", err)
		}

		ovfData = env.Properties[OVF_USERDATA]
		ovfEncoding = env.Properties[OVF_USERDATA_ENCODING]
	}

	data, err := info.String(GUESTINFO_USERDATA, ovfData)
	if err != nil {
		f.Logger.Debug("failed to fetch config: %v", err)
		return config{}, err
	}

	encoding, err := info.String(GUESTINFO_USERDATA_ENCODING, ovfEncoding)
	if err != nil {
		f.Logger.Debug("failed to fetch config encoding: %v", err)
		return config{}, err
	}

	return config{
		data:     data,
		encoding: encoding,
	}, nil
}

func DelConfig(f *resource.Fetcher) error {
	info := rpcvmx.NewConfig()

	// delete userdata if set and not already a deletion marker
	orig, err := info.String(GUESTINFO_USERDATA, GUESTINFO_DELETED_USERDATA)
	if err != nil {
		return fmt.Errorf("getting config property: %w", err)
	}
	if orig != GUESTINFO_DELETED_USERDATA {
		// we can't delete properties or set them to the empty
		// string, so set encoding to "base64" and data to encoded "{}"
		f.Logger.Info("deleting config from guestinfo properties")
		if err := info.SetString(GUESTINFO_USERDATA, GUESTINFO_DELETED_USERDATA); err != nil {
			return fmt.Errorf("replacing config: %w", err)
		}

		// overwrite encoding if unset or not already base64
		origEncoding, err := info.String(GUESTINFO_USERDATA_ENCODING, "")
		if err != nil {
			return fmt.Errorf("getting config encoding property: %w", err)
		}
		if origEncoding != GUESTINFO_DELETED_USERDATA_ENCODING {
			if err := info.SetString(GUESTINFO_USERDATA_ENCODING, GUESTINFO_DELETED_USERDATA_ENCODING); err != nil {
				return fmt.Errorf("replacing config encoding: %w", err)
			}
		}
	}

	ovfEnv, err := info.String(GUESTINFO_OVF, "")
	if err != nil {
		// unlike FetchConfig, don't ignore errors, since that could
		// have security implications
		return fmt.Errorf("reading OVF environment: %w", err)
	}
	if ovfEnv != "" {
		prunedData, didPrune, err := DeleteOvfProperties([]byte(ovfEnv), []string{OVF_USERDATA, OVF_USERDATA_ENCODING})
		if err != nil {
			return fmt.Errorf("deleting OVF properties: %w", err)
		}
		// don't rewrite the property if there's nothing to change
		if didPrune {
			f.Logger.Info("deleting config from OVF environment")
			if err := info.SetString(GUESTINFO_OVF, string(prunedData)); err != nil {
				return fmt.Errorf("replacing OVF environment: %w", err)
			}
		}
	}

	return nil
}
