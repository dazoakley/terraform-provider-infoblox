package infoblox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/infobloxopen/infoblox-go-client"
	"testing"
)

func TestAccResourceSRVRecord(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSRVRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceSRVRecordCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccSRVRecordExists(t, "infoblox_ptr_record.foo", "10.0.0.0/24", "10.0.0.2", "test", "demo-network", "default", "a.com"),
				),
			},
			resource.TestStep{
				Config: testAccresourceSRVRecordUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccSRVRecordExists(t, "infoblox_ptr_record.fcccoo", "10.0.0.0/24", "10.0.0.2", "test", "demo-network", "default", "a.com"),
				),
			},
		},
	})
}

func testAccCheckSRVRecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource_a_record" {
			continue
		}
		Connector := meta.(*ibclient.Connector)
		objMgr := ibclient.NewObjectManager(Connector, "terraform_test", "test")
		recordName, _ := objMgr.GetSRVRecordByRef(rs.Primary.ID)
		if recordName != nil {
			return fmt.Errorf("record not found")
		}

	}
	return nil
}
func testAccSRVRecordExists(t *testing.T, n string, cidr string, ipAddr string, networkViewName string, recordName string, dnsView string, zone string) resource.TestCheckFunc {
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

		recordName, _ := objMgr.GetSRVRecordByRef(rs.Primary.ID)
		if recordName == nil {
			return fmt.Errorf("record not found")
		}

		return nil
	}
}

var testAccresourceSRVRecordCreate = fmt.Sprintf(`
resource "infoblox_srv_record" "test_record"{
	  record_name = "demo1"
	  port = 17357
	  priority = 42
	  target = "something.demo"
	  dns_view="default"
	  tenant_id="test"
	  weight = 10
	  zone="aa.com"
	}`)

var testAccresourceSRVRecordUpdate = fmt.Sprintf(`
resource "infoblox_srv_record" "test_record"{
	  record_name = "demo1"
	  port = 17357
	  priority = 42
	  target = "something.demo"
	  dns_view="default"
	  tenant_id="test"
	  weight = 10
	  zone="aa.com"
	}`)
