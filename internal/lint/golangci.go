package lint

import (
	"fmt"

	"github.com/stalwartgiraffe/cmr/internal/config"
	"github.com/stalwartgiraffe/cmr/internal/xr"
)

func RunEach(cfg *config.Config) error {
	const expectedStatus int = 1
	for _, p := range cfg.Projects {
		for _, linter := range p.Linters {
			fmt.Println(">lint", linter.Name)
			out, err := xr.Run(linter.Name, expectedStatus, nil, linter.CmdArgs...)
			if err != nil {
				return err
			}
			fmt.Println(out)
		}
	}
	return nil
}
