# Artifact ID vault packer plugin

This is packer plugin to store artifact ID.

## Installation

### Using pre-built releases

#### Using `packer init`

To install this plugin, copy following code into your `.pkr.hcl` file.

```packer
packer {
  required_plguins {
    artifactidvault = {
      version = ">= 1.0.0"
      source = "github.com/cloudcore-tu/artifactidvault"
    }
  }
}
```

#### Manual installation

Pre-built binary can be found [here](https://github.com/cloudcore-tu/packer-plugin-artifactidvault/releases), so extract it to `~/.packer.d` as follows:

```bash
unzip "<archived pre-built binary path>"
cp "<extracted pre-built binary path>" ~/.packer.d
```

For more information, refer to [official documentation](https://www.packer.io/docs/plugins#installing-plugins).

### From sources

If you want to build from sources, clone GitHub repository and run command `go build` at project root.
Built successful, binary `packer-plugin-artifactidvault` is generated in project root.
You can then follow same procedure as `Manual installation` to install plugin.
