package testutil

import (
	"bytes"
	"io"

	"github.com/spf13/cobra"
)

type CommandOutput struct {
	Out string
	Err string
}

func ExecuteCommand(cmd *cobra.Command, args ...string) (CommandOutput, error) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return CommandOutput{
		Out: outBuf.String(),
		Err: errBuf.String(),
	}, err
}

func SetCommandIO(cmd *cobra.Command, out, errOut io.Writer) {
	cmd.SetOut(out)
	cmd.SetErr(errOut)
}
