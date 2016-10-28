package azurerm

import (
	"fmt"
	"log"

	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/redis"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmRedis() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRedisCreate,
		Read:   resourceArmRedisRead,
		Update: resourceArmRedisCreate,
		Delete: resourceArmRedisDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"redis_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"family": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRedisFamily,
			},

			"sku_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRedisSku,
			},

			"shard_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"enable_non_ssl_port": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"redis_configuration": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ssl_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmRedisCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisClient
	log.Printf("[INFO] preparing arguments for Azure ARM Redis creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)

	redisVersion := d.Get("redis_version").(string)
	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := redis.CreateOrUpdateParameters{
		Name:     &name,
		Location: &location,
		Properties: &redis.Properties{
			EnableNonSslPort: &enableNonSSLPort,
			RedisVersion:     &redisVersion,
			Sku: &redis.Sku{
				Capacity: &capacity,
				Family:   family,
				Name:     sku,
			},
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("shard_count"); ok {
		shardCount := int32(v.(int))
		parameters.Properties.ShardCount = &shardCount
	}

	/*
		if v, ok := d.GetOk("redis_configuration"); ok {
			params := v.(map[string]interface{})

			redisConfiguration := make(map[string]*string, len(params))
			for key, val := range params {
				redisConfiguration[key] = struct {
					Value *string
				}{
					Value: val.(string),
				}
			}

			parameters.Properties.RedisConfiguration = &redisConfiguration
		}
	*/

	_, err := client.CreateOrUpdate(resGroup, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Redis Instance %s (resource group %s) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for Redis Instance (%s) to become available", d.Get("name"))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(client, resGroup, name),
		Timeout:    30 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Redis Instance (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	return resourceArmRedisRead(d, meta)
}

func resourceArmRedisRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["Redis"]

	resp, err := client.Get(resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Redis %s: %s", name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	keysResp, err := client.ListKeys(resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making ListKeys request on Azure Redis %s: %s", name, err)
	}

	//log.Printf("[INFO] -----======-----")
	//log.Printf("[INFO] API: %s", spew.Sdump(resp))
	//log.Printf("[INFO] -----======-----")

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	parseAzureRMRedisProperties(d, resp.Properties)

	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	// TODO: Redis Configuation

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmRedisDelete(d *schema.ResourceData, meta interface{}) error {
	redisClient := meta.(*ArmClient).redisClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["Redis"]

	resp, err := redisClient.Delete(resGroup, name)

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("Error issuing Azure ARM delete request of Redis Instance '%s': %s", name, err)
	}

	return nil
}

func redisStateRefreshFunc(client redis.Client, resourceGroupName string, sgName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(resourceGroupName, sgName)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in redisStateRefreshFunc to Azure ARM for Redis Instance '%s' (RG: '%s'): %s", sgName, resourceGroupName, err)
		}

		return res, *res.Properties.ProvisioningState, nil
	}
}

func parseAzureRMRedisProperties(d *schema.ResourceData, properties *redis.ReadableProperties) {
	if properties != nil {
		d.Set("redis_version", properties.RedisVersion)
		d.Set("enable_non_ssl_port", properties.EnableNonSslPort)
		d.Set("ssl_port", properties.SslPort)

		d.Set("host_name", properties.HostName)

		if properties.Port != nil {
			d.Set("port", properties.Port)
		}

		if properties.Sku != nil {
			d.Set("capacity", properties.Sku.Capacity)
			d.Set("family", properties.Sku.Family)
			d.Set("sku_name", properties.Sku.Name)
		}

		// TODO: ensure this parses out correctly
		if properties.ShardCount != nil {
			d.Set("shard_count", properties.ShardCount)
		}
	}
}

func validateRedisFamily(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"c": true,
		"p": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Family can only be C or P"))
	}
	return
}

func validateRedisSku(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	skus := map[string]bool{
		"basic":    true,
		"standard": true,
		"premium":  true,
	}

	if !skus[value] {
		errors = append(errors, fmt.Errorf("Redis SKU can only be Basic, Standard or Premium"))
	}
	return
}
