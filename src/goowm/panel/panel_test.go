package panel_test

import (
	"goowm/panel"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockDisplayServer struct {
}

var _ = Describe("Panel", func() {
	It("returns a new Panel", func() {
		p := panel.New(&MockDisplayServer{})

		Expect(p).ToNot(BeNil())
	})
})
