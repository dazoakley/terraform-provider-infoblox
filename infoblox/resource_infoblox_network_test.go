package infoblox

import (
	"fmt"
	"testing"

	ibclient "github.com/alanplatt/infoblox-go-client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccresourceNetwork(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceNetworkCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateNetworkExists(t, "infoblox_network.foo", "10.10.0.0/24", "default", "demo-network"),
				),
			},
			resource.TestStep{
				Config: testAccresourceNetworkUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateNetworkExists(t, "infoblox_network.foo", "10.10.0.0/24", "default", "demo-network"),
				),
			},
		},
	})
}

func TestAccresourceNetwork_Allocate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceNetworkAllocate,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateNetworkExists(t, "infoblox_network.foo0", "10.0.0.0/24", "default", "demo-network"),
					testAccCreateNetworkExists(t, "infoblox_network.foo1", "10.0.1.0/24", "default", "demo-network"),
				),
			},
		},
	})
}

func testAccCheckNetworkDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "infoblox_network" {
			continue
		}
		Connector := meta.(*ibclient.Connector)
		objMgr := ibclient.NewObjectManager(Connector, "terraform_test", "test")
		networkName, _ := objMgr.GetNetwork("demo-network", "10.10.0.0/24", nil)
		if networkName != nil {
			return fmt.Errorf("Network not found")
		}

	}
	return nil
}

func testAccCreateNetworkExists(t *testing.T, n string, cidr string, networkViewName string, networkName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found:%s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID i set")
		}
		meta := testAccProvider.Meta()
		Connector := meta.(*ibclient.Connector)
		objMgr := ibclient.NewObjectManager(Connector, "terraform_test", "test")

		networkName, _ := objMgr.GetNetwork(networkName, cidr, nil)
		if networkName != nil {
			return fmt.Errorf("Network not found")
		}
		return nil
	}
}

var testAccresourceNetworkCreate = fmt.Sprintf(`
resource "infoblox_network" "foo"{
	network_view_name="default"
	network_name="demo-network"
	cidr="10.10.0.0/24"
	tenant_id="foo"
	}`)

/*
Right now no infoblox_network_container resource available
So, before run acceptance test TestAccresourceNetwork_Allocate
in default network view should be created network container 10.0.0.0/16
*/
var testAccresourceNetworkAllocate = fmt.Sprintf(`
resource "infoblox_network" "foo0"{
	network_view_name="default"
	network_name="demo-network"
	cidr="10.0.0.0/16"
	tenant_id="foo"
	allocate_prefix_len=24
	}
resource "infoblox_network" "foo1"{
	network_view_name="default"
	network_name="demo-network"
	cidr="10.0.0.0/16"
	tenant_id="foo"
	allocate_prefix_len=24
	}`)

var testAccresourceNetworkUpdate = fmt.Sprintf(`
resource "infoblox_network" "foo"{
	network_view_name="default"
	network_name="demo-network"
	cidr="10.10.0.0/24"
	tenant_id="foo"
	}`)
