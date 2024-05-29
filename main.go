/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>

This file is part of smarter-cli. smarter-cli is free software: you can
redistribute it and/or modify it under the terms of the MIT License.
smarter-cli is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE.  See the MIT License for more details.
*/
package main

import (
	"github.com/joho/godotenv"
	"github.com/smarter-sh/smarter-cli/cmd"
	_ "github.com/smarter-sh/smarter-cli/cmd/chat"
	_ "github.com/smarter-sh/smarter-cli/cmd/delete"
	_ "github.com/smarter-sh/smarter-cli/cmd/deploy"
	_ "github.com/smarter-sh/smarter-cli/cmd/describe"
	_ "github.com/smarter-sh/smarter-cli/cmd/get"
	_ "github.com/smarter-sh/smarter-cli/cmd/logs"
	_ "github.com/smarter-sh/smarter-cli/cmd/manifest"
	_ "github.com/smarter-sh/smarter-cli/cmd/undeploy"
)

var Version = "local.dev"

func main() {
	_ = godotenv.Load()

	cmd.Execute(Version)

}
