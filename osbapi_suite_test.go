package osbapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOsbapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OSBAPI Suite")
}
