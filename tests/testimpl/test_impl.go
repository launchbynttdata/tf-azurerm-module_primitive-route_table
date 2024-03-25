package common

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/lib/azure/configure"
	"github.com/launchbynttdata/lcaf-component-terratest/lib/azure/login"
	internalNetwork "github.com/launchbynttdata/lcaf-component-terratest/lib/azure/network"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

const terraformDir string = "../../examples/route_table"
const varFile string = "test.tfvars"

func TestRouteTable(t *testing.T, ctx types.TestContext) {

	envVarMap := login.GetEnvironmentVariables()
	clientID := envVarMap["clientID"]
	clientSecret := envVarMap["clientSecret"]
	tenantID := envVarMap["tenantID"]
	subscriptionID := envVarMap["subscriptionID"]

	spt, err := login.GetServicePrincipalToken(clientID, clientSecret, tenantID)
	if err != nil {
		t.Fatalf("Error getting Service Principal Token: %v", err)
	}

	routeTableClient := internalNetwork.GetRouteTablesClient(spt, subscriptionID)
	terraformOptions := configure.ConfigureTerraform(terraformDir, []string{terraformDir + "/" + varFile})
	t.Run("doesRouteTableExist", func(t *testing.T) {
		checkrouteTablesExistence(t, routeTableClient, terraformOptions, ctx)
	})
}

func checkrouteTablesExistence(t *testing.T, routeTableClient network.RouteTablesClient, terraformOptions *terraform.Options, ctx types.TestContext) {
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	routeTableName := terraform.Output(t, ctx.TerratestTerraformOptions(), "name")

	routeTable, err := routeTableClient.Get(context.Background(), resourceGroupName, routeTableName, "")
	if err != nil {
		t.Fatalf("Error getting Route Table: %v", err)
	}
	if routeTable.Name == nil {
		t.Fatalf("Route Table does not exist")
	}

	assert.Equal(t, *routeTable.Name, routeTableName)
}
