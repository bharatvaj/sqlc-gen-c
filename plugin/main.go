package main

import (
	"github.com/sqlc-dev/plugin-sdk-go/codegen"

	c "github.com/timwmillard/sqlc-gen-c/internal"
)

func main() {
	codegen.Run(c.Generate)
}
