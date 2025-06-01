package config

import (
	"embed"
	"io"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stalwartgiraffe/cmr/internal/utils"
)

//go:embed data/loadtests/*
var loadTestsFS embed.FS

var _ = Describe("for each data file, loadtests", func() {
	It("test every file", func() {
		err := utils.WalkFileReaders(loadTestsFS, func(path string, file io.Reader) {
			isOK := strings.HasPrefix(path, "data/loadtests/ok")
			cfg, err := LoadConfig(file)
			Expect(err == nil && cfg != nil).To(Equal(isOK))
		})
		Expect(err).To(BeNil())
	})
})
