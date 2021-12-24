package klustair

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func RunCli(ctx *cli.Context) error {
	fmt.Println("run")
	return nil
	//return xerrors.Errorf("option error: %w", "nothing to do")
}
