package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/goamz/ec2"
)

func resourceAwsVpcPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsVpcPeeringCreate,
		Read:   resourceAwsVpcPeeringRead,
		Update: resourceAwsVpcPeeringUpdate,
		Delete: resourceAwsVpcPeeringDelete,

		Schema: map[string]*schema.Schema{
			"peer_owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_accept": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"accept_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAwsVpcPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn

	// Create the vpc peering connection
	createOpts := &ec2.CreateVpcPeeringConnection{
		PeerOwnerId: d.Get("peer_owner_id").(string),
		PeerVpcId:   d.Get("peer_vpc_id").(string),
		VpcId:       d.Get("vpc_id").(string),
	}
	log.Printf("[DEBUG] VpcPeeringCreate create config: %#v", createOpts)
	resp, err := ec2conn.CreateVpcPeeringConnection(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating vpc peering connection: %s", err)
	}

	// Get the ID and store it
	rt := &resp.VpcPeeringConnection
	d.SetId(rt.VpcPeeringConnectionId)
	log.Printf("[INFO] Vpc Peering Connection ID: %s", d.Id())

	// Wait for the vpc peering connection to become available
	log.Printf(
		"[DEBUG] Waiting for vpc peering connection (%s) to become available",
		d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  "ready",
		Refresh: resourceAwsVpcPeeringConnectionStateRefreshFunc(ec2conn, d.Id()),
		Timeout: 1 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for vpc peering (%s) to become available: %s",
			d.Id(), err)
	}

	return resourceAwsVpcPeeringUpdate(d, meta)
}

func resourceAwsVpcPeeringRead(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn
	pcRaw, _, err := resourceAwsVpcPeeringConnectionStateRefreshFunc(ec2conn, d.Id())()
	if err != nil {
		return err
	}
	if pcRaw == nil {
		d.SetId("")
		return nil
	}

	pc := pcRaw.(*ec2.VpcPeeringConnection)

	if d.Get("auto_accept").(bool) {
		resourceVpcPeeringConnectionAccept(ec2conn, pc, d)
	} else {
		d.Set("accept_status", pc.Status.Code)
	}

	d.Set("peer_owner_id", pc.AccepterVpcInfo.OwnerId)
	d.Set("peer_vpc_id", pc.AccepterVpcInfo.VpcId)
	d.Set("vpc_id", pc.RequesterVpcInfo.VpcId)
	d.Set("tags", tagsToMap(pc.Tags))

	return nil
}

func resourceVpcPeeringConnectionAccept(conn *ec2.EC2, oldPc *ec2.VpcPeeringConnection, d *schema.ResourceData) error {
	if oldPc.Status.Code == "pending-acceptance" {
		log.Printf("[INFO] Accept Vpc Peering Connection with id: %s", d.Id())
		_, err := conn.AcceptVpcPeeringConnection(d.Id())
		if err != nil {
			return fmt.Errorf("Error accepting vpc peering connection: %s", err)
		}

		pcRaw, _, err := resourceAwsVpcPeeringConnectionStateRefreshFunc(conn, d.Id())()
		if err != nil {
			return err
		}
		if pcRaw == nil {
			d.SetId("")
			return nil
		}

		pc := pcRaw.(*ec2.VpcPeeringConnection)
		d.Set("accept_status", pc.Status.Code)

	}

	return nil
}

func resourceAwsVpcPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn

	if err := setTags(ec2conn, d); err != nil {
		return err
	} else {
		d.SetPartial("tags")
	}

	return resourceAwsVpcPeeringRead(d, meta)
}

func resourceAwsVpcPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn

	_, err := ec2conn.DeleteVpcPeeringConnection(d.Id())
	return err
}

// resourceAwsVpcPeeringConnectionStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a VpcPeeringConnection.
func resourceAwsVpcPeeringConnectionStateRefreshFunc(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		resp, err := conn.DescribeVpcPeeringConnection([]string{id}, ec2.NewFilter())
		if err != nil {
			if ec2err, ok := err.(*ec2.Error); ok && ec2err.Code == "InvalidVpcPeeringConnectionID.NotFound" {
				resp = nil
			} else {
				log.Printf("Error on VpcPeeringConnectionStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		pc := &resp.VpcPeeringConnections[0]

		return pc, "ready", nil
	}
}
