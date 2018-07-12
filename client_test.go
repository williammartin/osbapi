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

		AfterEach(func() {
			Expect(client.Deprovision("my-instance", &osbapi.DeprovisionRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})).To(Succeed())
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

	Describe("Creating a binding", func() {
		var (
			serviceID string
			planID    string
		)

		BeforeEach(func() {
			catalog, err := client.Catalog()
			Expect(err).NotTo(HaveOccurred())

			serviceID = catalog.Services[0].ID
			planID = catalog.Services[0].Plans[0].ID

			_, err = client.Provision("my-instance", &osbapi.ProvisionRequest{
				ServiceID:      serviceID,
				PlanID:         planID,
				OrganizationID: "org-id",
				SpaceID:        "space-id",
			})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(client.Unbind("my-instance", "my-binding", &osbapi.UnbindingRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})).To(Succeed())

			Expect(client.Deprovision("my-instance", &osbapi.DeprovisionRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})).To(Succeed())
		})

		It("results in a service binding", func() {
			_, err := client.Bind("my-instance", "my-binding", &osbapi.BindingRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})
			Expect(err).NotTo(HaveOccurred())

			binding, err := client.GetBinding("my-instance", "my-binding")
			Expect(err).NotTo(HaveOccurred())

			credentials, ok := binding.Credentials.(map[string]interface{})
			Expect(ok).To(BeTrue())

			Expect(credentials["username"]).NotTo(BeEmpty())
			Expect(credentials["password"]).NotTo(BeEmpty())
		})
	})

	Describe("Deleting a binding", func() {
		var (
			serviceID string
			planID    string
		)

		BeforeEach(func() {
			catalog, err := client.Catalog()
			Expect(err).NotTo(HaveOccurred())

			serviceID = catalog.Services[0].ID
			planID = catalog.Services[0].Plans[0].ID

			_, err = client.Provision("my-instance", &osbapi.ProvisionRequest{
				ServiceID:      serviceID,
				PlanID:         planID,
				OrganizationID: "org-id",
				SpaceID:        "space-id",
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = client.Bind("my-instance", "my-binding", &osbapi.BindingRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(client.Deprovision("my-instance", &osbapi.DeprovisionRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})).To(Succeed())
		})

		It("results in the service binding disappearing", func() {
			Expect(client.Unbind("my-instance", "my-binding", &osbapi.UnbindingRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})).To(Succeed())

			_, err := client.GetBinding("my-instance", "my-binding")
			Expect(err).To(MatchError(ContainSubstring("Not Found")))
		})
	})

	Describe("Deleting an instance", func() {
		var (
			serviceID string
			planID    string
		)

		BeforeEach(func() {
			catalog, err := client.Catalog()
			Expect(err).NotTo(HaveOccurred())

			serviceID = catalog.Services[0].ID
			planID = catalog.Services[0].Plans[0].ID

			_, err = client.Provision("my-instance", &osbapi.ProvisionRequest{
				ServiceID:      serviceID,
				PlanID:         planID,
				OrganizationID: "org-id",
				SpaceID:        "space-id",
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = client.Bind("my-instance", "my-binding", &osbapi.BindingRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(client.Unbind("my-instance", "my-binding", &osbapi.UnbindingRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})).To(Succeed())
		})

		It("results in the service instance disappearing", func() {
			Expect(client.Deprovision("my-instance", &osbapi.DeprovisionRequest{
				ServiceID: serviceID,
				PlanID:    planID,
			})).To(Succeed())

			_, err := client.GetInstance("my-instance")
			Expect(err).To(MatchError(ContainSubstring("Not Found")))
		})
	})
})
