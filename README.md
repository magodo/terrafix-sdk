# terrafix-sdk

A SDK to easier implement the supporting functions to a Terraform provider, which are required by the [`terrafix`](https://github.com/magodo/terrafix).

The main types are the following, which implement the `function.Function` interface defined in `"github.com/hashicorp/terraform-plugin-framework/function`:

- `function.UpgradeConfigDefinitionFunction`: This is the function that upgrades a Terraform configuration definition, for a provider, resource or data source.
- `function.UpgradeConfigReferenceFunction`: This is the function that upgrades Terraform configuration reference origins, targeting to a provider, resource or data source.

Check out the framework [document](https://developer.hashicorp.com/terraform/plugin/framework/functions/implementation) about how to implement these provider functions.
