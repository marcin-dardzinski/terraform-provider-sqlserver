package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/internal"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/sql"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider { return internal.Provider() },
		Debug:        debug,
	}

	plugin.Serve(opts)

	sql.DisposeConnections()
}
