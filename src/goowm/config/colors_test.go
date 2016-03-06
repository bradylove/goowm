package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"goowm/config"
)

var _ = Describe("Colors", func() {
	Describe("ParseHexColor", func() {
		It("can parse white", func() {
			res := config.ParseHexColor("#ffffff")
			Expect(res).To(Equal(0xffffff))
		})

		It("can parse black", func() {
			res := config.ParseHexColor("#000000")
			Expect(res).To(Equal(0x000000))
		})

		It("can parse red", func() {
			res := config.ParseHexColor("#ff3300")
			Expect(res).To(Equal(0xff3300))
		})
	})
})
