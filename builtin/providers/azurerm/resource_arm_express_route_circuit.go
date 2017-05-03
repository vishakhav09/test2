package azurerm

import (
	"log"

	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmExpressRouteCircuit() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmExpressRouteCircuitCreate,
		Read:   resourceArmExpressRouteCircuitRead,
		Update: resourceArmExpressRouteCircuitCreate,
		Delete: resourceArmExpressRouteCircuitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"service_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"peering_location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"bandwidth_in_mbps": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"sku_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(network.ExpressRouteCircuitSkuTierStandard),
				ValidateFunc: validateSkuTier,
			},

			"sku_family": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(network.MeteredData),
				ValidateFunc: validateSkuFamily,
			},

			"allow_classic_operations": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_provider_provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmExpressRouteCircuitCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	ercClient := client.expressRouteCircuitClient

	log.Printf("[INFO] preparing arguments for Azure ARM ExpressRouteCircuit creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	serviceProviderName := d.Get("service_provider_name").(string)
	peeringLocation := d.Get("peering_location").(string)
	bandwidthInMbps := int32(d.Get("bandwidth_in_mbps").(int))
	skuTier := d.Get("sku_tier").(string)
	skuFamily := d.Get("sku_family").(string)
	skuName := fmt.Sprintf("%s_%s", skuTier, skuFamily)
	allowRdfeOps := d.Get("allow_classic_operations").(bool)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	erc := network.ExpressRouteCircuit{
		Name:     &name,
		Location: &location,
		Sku: &network.ExpressRouteCircuitSku{
			Name:   &skuName,
			Tier:   network.ExpressRouteCircuitSkuTier(skuTier),
			Family: network.ExpressRouteCircuitSkuFamily(skuFamily),
		},
		ExpressRouteCircuitPropertiesFormat: &network.ExpressRouteCircuitPropertiesFormat{
			AllowClassicOperations: &allowRdfeOps,
			ServiceProviderProperties: &network.ExpressRouteCircuitServiceProviderProperties{
				ServiceProviderName: &serviceProviderName,
				PeeringLocation:     &peeringLocation,
				BandwidthInMbps:     &bandwidthInMbps,
			},
		},
		Tags: expandedTags,
	}

	_, err := ercClient.CreateOrUpdate(resGroup, name, erc, make(chan struct{}))
	if err != nil {
		return errwrap.Wrapf("Error Creating/Updating ExpressRouteCircuit {{err}}", err)
	}

	read, err := ercClient.Get(resGroup, name)
	if err != nil {
		return errwrap.Wrapf("Error Getting ExpressRouteCircuit {{err}}", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ExpressRouteCircuit %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmExpressRouteCircuitRead(d, meta)
}

func resourceArmExpressRouteCircuitRead(d *schema.ResourceData, meta interface{}) error {
	erc, resGroup, exists, err := retrieveErcByResourceId(d.Id(), meta)
	if err != nil {
		return err
	}

	if !exists {
		d.SetId("")
		log.Printf("[INFO] Express Route Circuit %q not found. Removing from state", d.Get("name").(string))
		return nil
	}

	d.Set("name", erc.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", erc.Location)

	if erc.ServiceProviderProperties != nil {
		d.Set("service_provider_name", erc.ServiceProviderProperties.ServiceProviderName)
		d.Set("peering_location", erc.ServiceProviderProperties.PeeringLocation)
		d.Set("bandwidth_in_mbps", erc.ServiceProviderProperties.BandwidthInMbps)
	}

	if erc.Sku != nil {
		d.Set("sku_tier", string(erc.Sku.Tier))
		d.Set("sku_family", string(erc.Sku.Family))
	}

	d.Set("service_provider_provisioning_state", string(erc.ServiceProviderProvisioningState))
	d.Set("service_key", erc.ServiceKey)
	d.Set("allow_classic_operations", erc.AllowClassicOperations)

	flattenAndSetTags(d, erc.Tags)

	return nil
}

func resourceArmExpressRouteCircuitDelete(d *schema.ResourceData, meta interface{}) error {
	ercClient := meta.(*ArmClient).expressRouteCircuitClient

	resGroup, name, err := extractResourceGroupAndErcName(d.Id())
	if err != nil {
		return errwrap.Wrapf("Error Parsing Azure Resource ID {{err}}", err)
	}

	_, err = ercClient.Delete(resGroup, name, make(chan struct{}))
	if err != nil {
		return errwrap.Wrapf("Error Deleting ExpressRouteCircuit {{err}}", err)
	}

	d.SetId("")
	return nil
}
