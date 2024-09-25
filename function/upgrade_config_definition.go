package function

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

type DefinitionUpgradeFunction func(version int, sbody *hclsyntax.Body, wbody *hclwrite.Body) error

type DefinitionUpgraders map[BlockType]map[string]DefinitionUpgradeFunction

type UpgradeConfigDefinitionFunction struct {
	Upgraders DefinitionUpgraders
}

var _ function.Function = UpgradeConfigDefinitionFunction{}

func NewUpgradeConfigDefinitionFunction(upgraders DefinitionUpgraders) function.Function {
	return &UpgradeConfigDefinitionFunction{Upgraders: upgraders}
}

func (a UpgradeConfigDefinitionFunction) Metadata(_ context.Context, _ function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "upgrade_config_definition"
}

func (a UpgradeConfigDefinitionFunction) Definition(_ context.Context, _ function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:             "Upgrade a Terraform config definition",
		Description:         "Upgrade a Terraform config definition for a provider, resource or data source",
		MarkdownDescription: "Upgrade a Terraform config definition for a provider, resource or data source",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "block_type",
				Description:         "Block type: provider, resource, datasource",
				MarkdownDescription: "Block type: provider, resource, datasource",
			},
			function.StringParameter{
				Name:                "block_name",
				Description:         "The block name (e.g. provider name, resource type)",
				MarkdownDescription: "The block name (e.g. provider name, resource type)",
			},
			function.Int64Parameter{
				Name:                "version",
				Description:         "The version of the schema, inferred from the Terraform state",
				MarkdownDescription: "The version of the schema, inferred from the Terraform state",
			},
			function.StringParameter{
				Name:                "raw_content",
				Description:         "The content of the block definition",
				MarkdownDescription: "The content of the block definition",
			},
		},
		Return: function.StringReturn{},
	}
}

func (a UpgradeConfigDefinitionFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var blockType, blockName string
	var version int64
	var rawContent string

	response.Error = function.ConcatFuncErrors(request.Arguments.Get(ctx, &blockType, &blockName, &version, &rawContent))
	if response.Error != nil {
		return
	}

	sf, diags := hclsyntax.ParseConfig([]byte(rawContent), "", hcl.InitialPos)
	if diags.HasErrors() {
		response.Error = function.NewFuncError(diags.Error())
		return
	}
	sbody := sf.Body.(*hclsyntax.Body).Blocks[0].Body

	wf, diags := hclwrite.ParseConfig([]byte(rawContent), "", hcl.InitialPos)
	if diags.HasErrors() {
		response.Error = function.NewFuncError(diags.Error())
		return
	}
	wbody := wf.Body().Blocks()[0].Body()

	var err error
	if m, ok := a.Upgraders[BlockType(blockType)]; ok {
		if u, ok := m[blockName]; ok {
			err = u(int(version), sbody, wbody)
		}
	}
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	response.Error = function.ConcatFuncErrors(response.Result.Set(ctx, string(wf.Bytes())))
}
