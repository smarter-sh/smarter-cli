/*
Copyright © 2024 Lawrence McDaniel lawrence@querium.com
*/
package main

import (
	"fmt"
	"os"

	"github.com/QueriumCorp/smarter-cli/cmd"
	_ "github.com/QueriumCorp/smarter-cli/cmd/delete"
	_ "github.com/QueriumCorp/smarter-cli/cmd/delete/account"
	_ "github.com/QueriumCorp/smarter-cli/cmd/get"
	_ "github.com/QueriumCorp/smarter-cli/cmd/get/account"
	_ "github.com/QueriumCorp/smarter-cli/cmd/manifest"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err == nil {
		fmt.Fprintln(os.Stderr, "Found and loaded .env file")
	}

	cmd.Execute()
}
