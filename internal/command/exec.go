package command

import (
	"github.com/hashibuto/artillery"
)

var ExecCommand = &artillery.Command{
	Name:        "exec",
	Description: "execute a shellcode payload directly",
	Arguments: []*artillery.Argument{
		{
			Name:        "payload",
			Description: "shellcode to be executed",
		},
	},
	OnExecute: execCmd,
}

func execCmd(n artillery.Namespace, p *artillery.Processor) error {
	var args struct {
		Payload string
	}
	err := artillery.Reflect(n, &args)
	if err != nil {
		return err
	}

	//syscall.Mmap(0)

	return nil
}
