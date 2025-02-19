package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPassGen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PassGen Suite")
}
