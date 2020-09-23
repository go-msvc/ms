# Micro-service Framework #

## A Basic Micro-Service
An example of a very basic micro-service is presented in examples/example1_1
It creates a domain with two operations and provide default domain config:
```
	domain := ms.NewDomain().
		WithConfig("greeter", config{Name: "Samuel"}).
		WithOper("hello", helloRequest{}).
		WithOper("cheers", cheersRequest{})
```
After that the domain is served with HTTP REST that will default to listen on localhost:12345:
```
	rest.New(domain).Run()
```
So the 'hello' service operation can be called from a terminal with:
```
curl -XGET http://localhost:12345/oper/hello
{"greeting":"Samuel says hello"}
```

## Domain ##
A domain is a collection of service operations.
You will note that the domain has no name in the code, and that is per design.
The idea is that a domain could be reusable for different purposes and the
name is defined during deployment. An example is a simple list service. The
same service can be deployed as a white-list or a black-list, and both will
have operations to add or remove entries, but the user of the list would not
need to know it is the same code providing both lists.

## Operation ##
Each operation has a name. In the above example: 'hello' and 'cheers'.
An operation is define by the request it consumes and the response it provides.
A request may be validated if it implements interface IValidator.

## Server ##
A domain is served with an IServer. The simple server provided is called 'rest',
but it is not a fully complient HTTP REST server - it breaks several rules.
It is only to demonstrates the concept and will change in future.
Another implementation of IServer is 'nats' to subscribe to a topic on NATS
and push responses back to the caller. Also just for demonstration at this point.

## Client ##
A service often has to call other services and need a client to provide the interface.
It may use different clients, depending on how the remote service is served.
The configuration of those clients should also be outside the scope of the operation
using them, so the IClient implementation for each remote service could be different
and change during and after deployment as the remote service is changed.
More details to follow...

## Config ##
Config may be loaded from different sources and the code should not care where
the values were obtained from. That should be left to the deployment.
Knowing the config that was used and where it came from is essential to maintain
a running service.
More details to follow...

## Logging ##
An essential part of every micro-service is how it logs what it does as this is the
only way to figure out if it works or does not work and if so, why not.

The logger is used for debug logs and should not be used in production.

Audit records may be written to capture the request and response structure. This can
be switched on/off during operation and be filtered to only write in some cases.

Metrics may be exposed for example to Prometheus which can indicate the path in
the code and increment counters when certain error conditions are hit. This is the
most efficient way of monitoring a otherwise working service.

More details to follow...