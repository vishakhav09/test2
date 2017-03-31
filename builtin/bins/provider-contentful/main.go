package main

import (
	"github.com/hashicorp/terraform/builtin/providers/contentful"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: contentful.Provider,
	})
}
