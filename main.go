/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
package main

import (
	"fmt"
	"os"

	"github.com/cloudcore/packer-plugin-artifactidvault/post-processor/ssm"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"
)

var (
	Version           = "1.0.0"
	VersionPrerelease = "dev"
	PluginVersion     = version.InitializePluginVersion(Version, VersionPrerelease)
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterPostProcessor("ssm", new(ssm.PostProcessor))
	pps.SetVersion(PluginVersion)
	if err := pps.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
