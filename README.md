### OSBAPI Client

This is a simple implementation of some subset of a client to talk to OSBAPI 
service brokers. It's a toy that I'm using to drive out my understanding of
the specification. Use at your peril.

## Running the tests

You can run the tests using `ginkgo`, but they won't run in parallel because
they will totally pollute each other. You need a clean instance of the
[overview-broker](https://github.com/mattmcneeney/overview-broker) running
on port 3000.
