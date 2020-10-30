package infoblox

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func resourceSRVRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceSRVRecordCreate,
		Read:   resourceSRVRecordGet,
		Update: resourceSRVRecordUpdate,
		Delete: resourceSRVRecordDelete,

		Schema: map[string]*schema.Schema{

			"record_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the record.",
			},

			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port of the service.",
			},

			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The Priority of the record.",
			},

			"target": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "The target of the record.",
			},

			"dns_view": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Dns View under which the zone has been created.",
			},

			"weight": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The weighting of the record.",
			},

			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of your tenant in cloud.",
			},

			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone under which record has to be created.",
			},
		},
	}
}

func resourceSRVRecordCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning to create A record from  required network block", resourceSRVRecordIDString(d))

	recordName := d.Get("record_name").(string)
	port := uint(d.Get("port").(int))
	priority := uint(d.Get("priority").(int))
	target := d.Get("target").(string)
	weight := uint(d.Get("weight").(int))
	tenantID := d.Get("tenant_id").(string)
	dnsView := d.Get("dns_view").(string)
	zone := d.Get("zone").(string)
	connector := m.(*ibclient.Connector)

	ea := make(ibclient.EA)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)
	name := recordName + "." + zone
	recordSRV, err := objMgr.CreateSRVRecord(dnsView, name, port, priority, target, weight, ea)

	if err != nil {
		return fmt.Errorf("Error creating SRV Record (%s): %s", name, err)
	}

	d.Set("recordName", name)
	d.SetId(recordSRV.Ref)

	log.Printf("[DEBUG] %s: Creation of A Record complete", resourceSRVRecordIDString(d))

	return resourceSRVRecordGet(d, m)
}

func resourceSRVRecordGet(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning to Get SRV Record", resourceSRVRecordIDString(d))

	dnsView := d.Get("dns_view").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.GetSRVRecordByRef(d.Id())
	if err != nil {
		return fmt.Errorf("Getting SRV record failed from dns view (%s) : %s", dnsView, err)
	}
	d.SetId(obj.Ref)
	log.Printf("[DEBUG] %s: Completed reading required SRV Record ", resourceSRVRecordIDString(d))
	return nil
}

func resourceSRVRecordUpdate(d *schema.ResourceData, m interface{}) error {

	// return fmt.Errorf("updating a SRV record is not supported")

	log.Printf("[DEBUG] %s: Beginning to Update SRV Record", resourceSRVRecordIDString(d))

	dnsView := d.Get("dns_view").(string)
	port := uint(d.Get("port").(int))
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.UpdateSRVRecord(d.Id(), dnsView, port)
	if err != nil {
		return fmt.Errorf("Updating SRV record failed from dns view (%s) : %s", dnsView, err)
	}
	d.SetId(obj.Ref)
	log.Printf("[DEBUG] %s: Completed updating required SRV Record ", resourceSRVRecordIDString(d))
	return nil
}

func resourceSRVRecordDelete(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] %s: Beginning Deletion of SRV Record", resourceSRVRecordIDString(d))

	dnsView := d.Get("dns_view").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	_, err := objMgr.DeleteSRVRecord(d.Id())
	if err != nil {
		return fmt.Errorf("Deletion of SRV Record failed from dns view(%s) : %s", dnsView, err)
	}
	d.SetId("")

	log.Printf("[DEBUG] %s: Deletion of SRV Record complete", resourceSRVRecordIDString(d))
	return nil
}

type resourceSRVRecordIDStringInterface interface {
	Id() string
}

func resourceSRVRecordIDString(d resourceSRVRecordIDStringInterface) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("infoblox_srv_record (ID = %s)", id)
}
