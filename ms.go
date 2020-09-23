package ms

import (
	"regexp"
)

type IServer interface {
	Run() error
}

type IOper interface {
	Run(ctx IContext) (IResult, IResponse)
}

type IResult interface {
	Error() error
}

type IResultCode interface {
	Code() int
	Desc() string
}

type IResponse interface{}

type IValidator interface {
	Validate() error
}

const namePattern = `[a-zA-Z]([A-Za-z0-9_-]*[a-zA-Z0-9])*`

var nameRegex = regexp.MustCompile("^" + namePattern + "$")

type IIdGen interface {
	New() string
}
