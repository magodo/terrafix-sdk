package function

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/terrafix-sdk/internal/hclutils"
)

type ReferenceUpgradeFunction func(version int, traversals []hcl.Traversal) ([]hcl.Traversal, error)

type ReferenceUpgraders map[BlockType]map[string]ReferenceUpgradeFunction

type UpgradeConfigReferenceFunction struct {
	Upgraders ReferenceUpgraders
}

var _ function.Function = UpgradeConfigReferenceFunction{}

func NewUpgradeConfigReferenceFunction(upgraders ReferenceUpgraders) function.Function {
	return &UpgradeConfigReferenceFunction{Upgraders: upgraders}
}

func (a UpgradeConfigReferenceFunction) Metadata(_ context.Context, _ function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "upgrade_config_references"
}

func (a UpgradeConfigReferenceFunction) Definition(_ context.Context, _ function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:             "Upgrade Terraform config reference origins",
		Description:         "Upgrade Terraform config reference origins targeting to a provider, resource or data source",
		MarkdownDescription: "Upgrade Terraform config reference origins targeting to a provider, resource or data source",
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
			function.ListParameter{
				Name:                "raw_contents",
				Description:         "The list of reference origin contents",
				MarkdownDescription: "The list of reference origin contents",
				ElementType:         basetypes.StringType{},
			},
		},
		Return: function.ListReturn{
			ElementType: basetypes.StringType{},
		},
	}
}

func (a UpgradeConfigReferenceFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var blockType, blockName string
	var version int
	var rawContents []string

	response.Error = function.ConcatFuncErrors(request.Arguments.Get(ctx, &blockType, &blockName, &version, &rawContents))
	if response.Error != nil {
		return
	}

	var traversals []hcl.Traversal
	for _, content := range rawContents {
		expr, diags := hclsyntax.ParseExpression([]byte(content), "", hcl.InitialPos)
		if diags.HasErrors() {
			response.Error = function.NewFuncError(diags.Error())
			return
		}
		var tv hcl.Traversal
		switch expr := expr.(type) {
		case *hclsyntax.ScopeTraversalExpr:
			tv = expr.AsTraversal()
		case *hclsyntax.RelativeTraversalExpr:
			tv = expr.AsTraversal()
		default:
			response.Error = function.NewFuncError(fmt.Sprintf("unexpected non-traversal expression: %s", content))
			return
		}
		traversals = append(traversals, tv)
	}

	if m, ok := a.Upgraders[BlockType(blockType)]; ok {
		if u, ok := m[blockName]; ok {
			var err error
			traversals, err = u(int(version), traversals)
			if err != nil {
				response.Error = function.NewFuncError(err.Error())
				return
			}
		}
	}

	var updateContents []string
	for _, tv := range traversals {
		updateContents = append(updateContents, hclutils.FormatTraversal(tv))
	}
	response.Error = function.ConcatFuncErrors(response.Result.Set(ctx, updateContents))
	return
}
