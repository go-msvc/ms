//Package main creates a simple microservice with two operations called 'hello' and 'cheers'
//You can build it simply with 'go build' and run with './example1_1' to start serving
//it with a HTTP REST interface that defaults to localhost:12345
//then request either operation with:
//  curl -XGET http://localhost:12345/oper/hello
//  curl -XGET http://localhost:12345/oper/cheers
package main

import (
	"fmt"

	"github.com/go-msvc/ms"
	"github.com/go-msvc/ms/rest"
)

func main() {
	//this is the minimal effort you need to define a micro-service
	domain := ms.New(nil).
		WithOper("hello", helloRequest{}).
		WithOper("cheers", cheersRequest{})

	//and serve it on a rest interface that defaults to address localhost:12345
	//next examples will show you how to configure things
	rest.New(domain).Run()
}

type helloRequest struct{}

type helloResponse struct {
	Greeting string `json:"greeting"`
}

func (helloRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	ctx.Debugf("hello")
	return nil, helloResponse{
		Greeting: fmt.Sprintf("Hello"),
	}
}

type cheersRequest struct{}

func (cheersRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	ctx.Debugf("cheers")
	//let operation fail like this:
	return ms.Result{Code: 1, Desc: "NYI", Details: "Cheers is not yet implemented"}, nil
}
