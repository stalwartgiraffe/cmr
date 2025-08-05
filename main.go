package main

import (
	"os"

	"github.com/stalwartgiraffe/cmr/internal/run"
)

func main() {
	os.Exit(run.NewRunner().Run())
}
