/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
//go:generate packer-sdc mapstructure-to-hcl2 -type Config
//go:generate packer-sdc struct-markdown

package ssm

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

const BuilderId = "packer.post-process.artifact-id-vault"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// SSM parameter name
	ParameterName string `mapstructure:"parameter-name" required:"true"`
	// Flag whether encryption is applied to the value.
	// If true, value type is SecureString, otherwise String.
	Encrypt bool `mapstructure:"encrypt"`
	// Region where value is put
	Region string `mapstructure:"region"`
	// Flag where to overwrite
	// for more information, see [PutParameter AWS API docs](https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_PutParameter.html#systemsmanager-PutParameter-request-Overwrite)
	Overwrite bool `mapstructure:"overwrite"`
	// Configuration value for extracting string from artifact ID.
	// use [regexp.FindString](https://pkg.go.dev/regexp#Regexp.FindString) to extract string that match specified regular expression.
	Matcher string `mapstructure:"matcher" required:"true"`
	ctx     interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	dcOpt := &config.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}
	if err := config.Decode(&p.config, dcOpt, raws...); err != nil {
		return err
	}

	if p.config.Matcher == "" {
		return fmt.Errorf("match must be non-empty string")
	}
	if p.config.ParameterName == "" {
		return fmt.Errorf("parameter-name must be non-empty string")
	}
	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, source packer.Artifact) (packer.Artifact, bool, bool, error) {
	opts := []func(*awsconfig.LoadOptions) error{}
	if p.config.Region != "" {
		opts = append(opts, awsconfig.WithRegion(p.config.Region))
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return source, false, false, err
	}
	svc := ssm.NewFromConfig(cfg)
	if err != nil {
		return source, false, false, err
	}
	typ := types.ParameterTypeString
	if p.config.Encrypt {
		typ = types.ParameterTypeSecureString
	}
	r, err := regexp.Compile(p.config.Matcher)
	if err != nil {
		return source, false, false, fmt.Errorf("matcher syntax error: %s", err.Error())
	}
	val := r.FindString(source.Id())
	param := &ssm.PutParameterInput{
		Name:      aws.String(p.config.ParameterName),
		Type:      typ,
		Value:     aws.String(val),
		Overwrite: aws.Bool(p.config.Overwrite),
	}
	if _, err := svc.PutParameter(ctx, param); err != nil {
		return source, false, false, err
	}
	return source, true, true, nil
}
