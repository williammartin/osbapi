package osbapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/williammartin/osbapi"
)

var _ = Describe("Client", func() {

	var (
		client *osbapi.Client
	)

	BeforeEach(func() {
		// start an instance of the overview broker
		client = osbapi.NewClient("http://127.0.0.1:3000",
			osbapi.WithAPIVersion("2.13"),
			osbapi.WithBasicAuth("admin", "password"))
	})

	Describe("Fetching the catalog", func() {

		It("can parse an OSBAPI catalog response", func() {
			catalog, err := client.Catalog()
			Expect(err).NotTo(HaveOccurred())
			Expect(catalog.Services[0].Name).To(Equal("overview-broker"))
			Expect(catalog.Services[0].Plans[0].Name).To(Equal("simple"))
		})
	})

	Describe("Provisioning a service", func() {
		var (
			serviceID string
			planID    string
		)

		BeforeEach(func() {
			catalog, err := client.Catalog()
			Expect(err).NotTo(HaveOccurred())

			serviceID = catalog.Services[0].ID
			planID = catalog.Services[0].Plans[0].ID
		})

		It("results in a service instance", func() {
			_, err := client.Provision("my-instance", &osbapi.ProvisionRequest{
				ServiceID:      serviceID,
				PlanID:         planID,
				OrganizationID: "org-id",
				SpaceID:        "space-id",
			})
			Expect(err).NotTo(HaveOccurred())

			instance, err := client.GetInstance("my-instance")
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.ServiceID).To(Equal(serviceID))
			Expect(instance.PlanID).To(Equal(planID))
		})
	})
})
