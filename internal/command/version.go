package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
)

var Version string

var VersionCommand = &artillery.Command{
	Name:        "version",
	Description: "display version information",
	OnExecute: func(n artillery.Namespace, p *artillery.Processor) error {
		fmt.Println(Version)
		return nil
	},
}
