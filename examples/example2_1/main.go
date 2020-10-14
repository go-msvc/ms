package main

import (
	"github.com/go-msvc/config"
	"github.com/go-msvc/config/source/mem"
	"github.com/go-msvc/logger"
	"github.com/go-msvc/ms"
	"github.com/go-msvc/ms/rest"
)

var (
	resultNotYetImplemented = ms.NewResultCode(1, "NOT_YET_IMPLEMENTED")
)

func main() {
	//hard coded config for example:
	config.Sources().Reset()
	cfgInMemory := mem.New().
		With("name", "example2_1").
		With("numbers", map[string]interface{}{"integers": []int{1, 2, 3}})
	config.Sources().Add(cfgInMemory)

	//log to terminal
	logger.Top().WithStream(logger.Terminal(logger.LogLevelDebug))

	//create a domain here passing in configurable struct myDomain{}
	domain := ms.New(myDomain{}).
		WithOper("sum", sumRequest{}).
		WithOper("prod", prodRequest{}).
		WithOper("div", divRequest{}).
		WithOper("name", nameRequest{})
	if err := rest.New(domain).Run(); err != nil {
		panic(err)
	}
}

//the domain struct is configurable and its value becomes accessible to each new context
//so the value won't change during execution
type myDomain struct {
	Name    string    `json:"name"`
	Numbers myNumbers `json:"numbers"`
}

//myNumbers is configurable
type myNumbers struct {
	sum  int
	prod int
}

//constructor from config
func (myNumbers) Create(cfg myNumbersConfig) (*myNumbers, error) {
	n := &myNumbers{
		sum:  0,
		prod: 1,
	}
	for _, i := range cfg.Integers {
		n.sum += i
		n.prod *= i
	}
	return n, nil
}

type myNumbersConfig struct {
	Integers []int `json:"integers"`
}

//=========================================================================
type sumRequest struct{}
type sumResponse struct {
	Total int `json:"total"`
}

func (req sumRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	n := ctx.DomainStruct().(myDomain)
	return nil, sumResponse{Total: n.Numbers.sum}
}

type prodRequest struct{}
type prodResponse struct {
	Total int `json:"total"`
}

func (req prodRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	n := ctx.DomainStruct().(myDomain)
	return nil, prodResponse{Total: n.Numbers.prod}
}

type divRequest struct{}

func (req divRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	//this oper always fail to illustrate:
	return ms.ErrorResult(
		nil,
		resultNotYetImplemented,
	), nil
}

type nameRequest struct{}
type nameResponse struct {
	Name string `json:"name"`
}

func (req nameRequest) Run(ctx ms.IContext) (ms.IResult, ms.IResponse) {
	n := ctx.DomainStruct().(myDomain)
	return nil, nameResponse{Name: n.Name}
}
