package command

// int runShellcode(void* buf) {
//   int (*funcptr)() = buf;
//   return funcptr();
// }
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/utils"
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

	b, err := utils.ParseHexArrayString(args.Payload)
	if err != nil {
		return err
	}

	targBuf, err := syscall.Mmap(-1, 0, len(b), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		return err
	}

	copy(targBuf, b)
	retCode := int(C.runShellcode(unsafe.Pointer(&targBuf)))
	fmt.Printf("Ret %d\n", retCode)

	return nil
}
