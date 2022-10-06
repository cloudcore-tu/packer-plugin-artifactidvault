package main

import (
	"fmt"
	"os"

	"github.com/cloudcore/packer-plugin-artifactidvault/ssm"
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