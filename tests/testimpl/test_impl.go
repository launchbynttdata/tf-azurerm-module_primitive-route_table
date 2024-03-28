package common

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestRouteTable(t *testing.T, ctx types.TestContext) {

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID is not set in the environment variables ")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %e\n", err)
	}

	clientFactory, err := armnetwork.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		t.Fatalf("Unable to get clientFactory: %e\n", err)
	}
	routeTableClient := clientFactory.NewRouteTablesClient()
	t.Run("doesRouteTableExist", func(t *testing.T) {
		checkRouteTablesExistence(t, routeTableClient, ctx)
	})
}

func checkRouteTablesExistence(t *testing.T, routeTableClient *armnetwork.RouteTablesClient, ctx types.TestContext) {
	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	routeTableName := terraform.Output(t, ctx.TerratestTerraformOptions(), "name")

	routeTable, err := routeTableClient.Get(context.Background(), resourceGroupName, routeTableName, nil)
	if err != nil {
		t.Fatalf("Error getting Route Table: %v", err)
	}
	if routeTable.Name == nil {
		t.Fatalf("Route Table does not exist")
	}

	assert.Equal(t, *routeTable.Name, routeTableName)
}
