/*
Copyright © 2024 Lawrence McDaniel lawrence@querium.com
*/
package main

import (
	"github.com/QueriumCorp/smarter-cli/cmd"
	_ "github.com/QueriumCorp/smarter-cli/cmd/delete"
	_ "github.com/QueriumCorp/smarter-cli/cmd/deploy"
	_ "github.com/QueriumCorp/smarter-cli/cmd/describe"
	_ "github.com/QueriumCorp/smarter-cli/cmd/get"
	_ "github.com/QueriumCorp/smarter-cli/cmd/logs"
	_ "github.com/QueriumCorp/smarter-cli/cmd/manifest"
	_ "github.com/QueriumCorp/smarter-cli/cmd/undeploy"
	"github.com/joho/godotenv"
)

var Version = "local.dev"

func main() {
	_ = godotenv.Load()

	cmd.Execute(Version)

}
