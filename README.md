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
