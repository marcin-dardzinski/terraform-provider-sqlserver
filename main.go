package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/sql"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider { return Provider() },
	}

	if debug {
		err := plugin.Debug(context.Background(), "local/local/sqlserver", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		plugin.Serve(opts)
	}

	sql.DisposeConnections()
}
