package ms

import (
	"regexp"
	"sync"

	"github.com/go-msvc/errors"
)

type Result struct {
	Code    int
	Desc    string
	Details string
}

func (r Result) Error() error {
	if r.Code%100 == 0 {
		return nil
	}
	return errors.Errorf("%s: %s", r.Desc, r.Details)
}

const descPattern = `^[A-Z0-9_]+$`

var (
	resultCodes = map[int]IResultCode{
		0: resultCode{0, "SUCCESS"},
	}
	resultCodeMutex sync.Mutex
	descRegex       = regexp.MustCompile(descPattern)
)

type resultCode struct {
	code int
	desc string
}

func (r resultCode) Code() int    { return r.code }
func (r resultCode) Desc() string { return r.desc }

func NewResultCode(code int, desc string) IResultCode {
	resultCodeMutex.Lock()
	defer resultCodeMutex.Unlock()
	newResultCode := resultCode{code: code, desc: desc}
	return addResultCode(newResultCode)
}

func addResultCode(r IResultCode) IResultCode {
	if !descRegex.MatchString(r.Desc()) {
		panic(errors.Errorf("result code %d desc \"%s\" may only contain [A-Z0-9_]", r.Code(), r.Desc()))
	}
	if existing, ok := resultCodes[r.Code()]; ok {
		if existing.Code() != r.Code() || existing.Desc() != r.Desc() {
			panic(errors.Errorf("duplicate ResultCode(%d:%s) already registered != new code (%d:%s)", existing.Code(), existing.Desc(), r.Code(), r.Desc()))
		}
		return existing
	}
	resultCodes[r.Code()] = r
	return r
}
