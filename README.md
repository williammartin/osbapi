### OSBAPI Client

This is a simple implementation of some subset of a client to talk to OSBAPI 
service brokers. It's a toy that I'm using to drive out my understanding of
the specification. Use at your peril.

## Running the tests

You can run the tests using `ginkgo`, but they won't run in parallel because
they will totally pollute each other. You need a clean instance of the
[overview-broker](https://github.com/mattmcneeney/overview-broker) running
on port 3000.

## Simple usage

The following provides a trimmed (no error checks) set of code to fetch the 
catalog, provision, bind, unbind, unprovision with some fetching of instances and
bindings.

```go
client := osbapi.NewClient("http://my-broker.com:8080,
    osbapi.WithAPIVersion("2.13"),
    osbapi.WithBasicAuth("admin", "password"))

catalog, _ := client.GetCatalog()

serviceID := catalog.Services[0].ID
planID := catalog.Services[0].Plans[0].ID

client.Provision("my-instance", &osbapi.ProvisionRequest{
    ServiceID:      serviceID,
    PlanID:         planID,
    OrganizationID: "org-id",
    SpaceID:        "space-id",
})

client.Bind("my-instance", "my-binding", &osbapi.BindingRequest{
    ServiceID: serviceID,
    PlanID:    planID,
})

client.GetInstance("my-instance")

client.GetBinding("my-instance", "my-binding")

client.Unbind("my-instance", "my-binding", &osbapi.UnbindingRequest{
    ServiceID: serviceID,
    PlanID:    planID,
})

client.Deprovision("my-instance", &osbapi.DeprovisionRequest{
    ServiceID: serviceID,
    PlanID:    planID,
})
```
