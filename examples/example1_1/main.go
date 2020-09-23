//Package main creates a simple microservice with two operations called 'hello' and 'cheers'
//You can build it simply with 'go build' and run with './example1_1' to start serving
//it with a HTTP REST interface that defaults to localhost:12345
//then request either operation with:
//  curl -XGET http://localhost:12345/oper/hello
//  curl -XGET http://localhost:12345/oper/cheers
package main

import (
	"fmt"

	"github.com/go-msvc/errors"
	"github.com/go-msvc/ms"
	"github.com/go-msvc/ms/rest"
)

func main() {
	domain := ms.NewDomain().
		WithConfig("greeter", config{Name: "Samuel"}).
		WithOper("hello", helloRequest{}).
		WithOper("cheers", cheersRequest{})

	rest.New(domain).Run()
}

type config struct {
	Name string `json:"name"`
}

func (c *config) Validate() error {
	if len(c.Name) == 0 {
		return errors.Errorf("missing name")
	}
	return nil
}

type helloRequest struct{}

type helloResponse struct {
	Greeting string `json:"greeting"`
}

func (o helloRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	ctx.Debugf("hello")
	cfg := ctx.GetConfig("greeter").(config)
	return nil, helloResponse{
		Greeting: fmt.Sprintf("%s says hello", cfg.Name),
	}
}

type cheersRequest struct{}

func (o cheersRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	ctx.Debugf("cheers")
	return ms.Result{Code: 1, Desc: "NYI"}, nil
}
