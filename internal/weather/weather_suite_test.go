package weather_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWeather(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Weather Suite")
}
