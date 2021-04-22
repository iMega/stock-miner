package acceptance

import (
	"errors"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/imega/stock-miner/tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var GraphQLUrl = "http://app/query"

var _ = BeforeSuite(func() {
	helpers.MockHTTPServer.BeforeSuite()

	err := WaitForSystemUnderTestReady()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	helpers.MockHTTPServer.AfterSuite()
})

func WaitForSystemUnderTestReady() error {
	req, err := http.NewRequest(http.MethodGet, "http://app/healthcheck", nil)
	if err != nil {
		return err
	}

	for attempts := 30; attempts > 0; attempts-- {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if err == nil && resp != nil && resp.StatusCode == http.StatusNoContent {
			return nil
		}

		log.Printf("ATTEMPTING TO CONNECT")

		<-time.After(time.Second)
	}

	return errors.New("SUT is not ready for tests")
}

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}
