package opsgenie

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk/user"
)

func resourceOpsGenieUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieUserCreate,
		Read:   resourceOpsGenieUserRead,
		Update: resourceOpsGenieUserUpdate,
		Delete: resourceOpsGenieUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				// TODO: validation as it's an email address, max length 100
			},
			"full_name": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: Max length 255
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: Max length 255
			},

			// TODO: Locale & Timezone
		},
	}
}

func resourceOpsGenieUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).users

	username := d.Get("username").(string)
	fullName := d.Get("full_name").(string)
	role := d.Get("role").(string)

	createRequest := user.CreateUserRequest{
		Username: username,
		Fullname: fullName,
		Role:     role,
	}

	log.Printf("[INFO] Creating OpsGenie user '%s'", username)

	createResponse, err := client.Create(createRequest)
	if err != nil {
		return err
	}

	err = checkOpsGenieResponse(createResponse.Code, createResponse.Status)
	if err != nil {
		return err
	}

	getRequest := user.GetUserRequest{
		Username: username,
	}

	getResponse, err := client.Get(getRequest)
	if err != nil {
		return err
	}

	d.SetId(getResponse.Id)

	return resourceOpsGenieUserRead(d, meta)
}

func resourceOpsGenieUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).users

	listRequest := user.ListUsersRequest{}
	listResponse, err := client.List(listRequest)
	if err != nil {
		return err
	}

	var found *user.GetUserResponse
	for _, user := range listResponse.Users {
		if user.Id == d.Id() {
			found = &user
			break
		}
	}

	if found == nil {
		d.SetId("")
		log.Printf("[INFO] User %q not found. Removing from state", d.Get("username").(string))
		return nil
	}

	getRequest := user.GetUserRequest{
		Id: d.Id(),
	}

	getResponse, err := client.Get(getRequest)
	if err != nil {
		return err
	}

	d.Set("username", getResponse.Username)
	d.Set("full_name", getResponse.Fullname)
	d.Set("role", getResponse.Role)

	return nil
}

func resourceOpsGenieUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).users

	updateRequest := user.UpdateUserRequest{
		Id: d.Id(),
	}

	username := d.Get("username").(string)
	fullName := d.Get("full_name").(string)
	role := d.Get("role").(string)

	if d.HasChange("full_name") {
		updateRequest.Fullname = fullName
	}

	if d.HasChange("role") {
		updateRequest.Role = role
	}

	log.Printf("[INFO] Updating OpsGenie user '%s'", username)

	updateResponse, err := client.Update(updateRequest)
	if err != nil {
		return err
	}

	err = checkOpsGenieResponse(updateResponse.Code, updateResponse.Status)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieUserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie user '%s'", d.Get("username").(string))
	client := meta.(*OpsGenieClient).users

	deleteRequest := user.DeleteUserRequest{
		Id: d.Id(),
	}

	_, err := client.Delete(deleteRequest)
	if err != nil {
		return err
	}

	return nil
}
