/*
Copyright © 2024 Lawrence McDaniel lawrence@querium.com
*/
package main

import (
	"github.com/QueriumCorp/smarter-cli/cmd"
	_ "github.com/QueriumCorp/smarter-cli/cmd/delete"
	_ "github.com/QueriumCorp/smarter-cli/cmd/delete/account"
	_ "github.com/QueriumCorp/smarter-cli/cmd/get"
	_ "github.com/QueriumCorp/smarter-cli/cmd/get/account"
)

func main() {
	cmd.Execute()
}
