# terrafix-sdk

A SDK to easier implement the supporting functions to a Terraform provider, which are required by the [`terrafix`](https://github.com/magodo/terrafix).

## Provider Functions

The main types are the followings, which implement the `function.Function` interface defined in `"github.com/hashicorp/terraform-plugin-framework/function`:

- `tfxsdk.NewFixConfigDefinitionFunction`: This returns the provider function `terrafix_config_definition` that fixes a Terraform configuration definition, for a provider, resource or data source.
- `tfxsdk.NewFixConfigReferenceFunction`: This returns the provider function `terrafix_config_references` that fixes Terraform configuration reference origins, targeting to a provider, resource or data source.

Check out the framework [document](https://developer.hashicorp.com/terraform/plugin/framework/functions/implementation) about how to register these provider functions.

## Helpers 

This module also provides some helpers to help users to implement configuration definition/reference fixers:

- *tfxsdk/traversal.go*: HCL Traversal related functions
