package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/go-msvc/errors"
	"github.com/go-msvc/logger"
	"github.com/go-msvc/ms"
)

//todo: unknown operation does not respond with NotFound!
//todo: ensure loggers are released when context ends
//todo: logger names not reflect package name? e.g. .../logger/domain/...
//todo: context must timeout with error result
//todo: configurable server address
//todo: configurable server type e.g. rest vs nats or http+html

func New(domain ms.IDomain) ms.IServer {
	return server{
		ILogger: domain.NewLogger("rest"),
		domain:  domain,
		idGen:   ms.NewIdGen(),
	}
}

//server implements ms.IServer
type server struct {
	logger.ILogger
	domain ms.IDomain
	idGen  ms.IIdGen
}

func (s server) Run() error {
	http.ListenAndServe("localhost:12345", s)
	return nil
}

func (s server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	id := s.idGen.New()
	ctx := context{
		IContext: s.domain.NewContext(id),
		res:      res,
		req:      req,
		path:     pathSplit(req.URL.Path),
	}
	defer ctx.Cancel()

	var err error
	switch ctx.path.Part(0) {
	case "oper":
		err = s.serveOper(ctx)
	default:
		err = errors.Errorf("path must start with /oper/<oper>")
	}
	if err != nil {
		ctx.Debugf("Serve Failed: %+v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	ctx.Debugf("Serve Success")
}

type context struct {
	ms.IContext
	res  http.ResponseWriter
	req  *http.Request
	path path
}

func (s server) serveOper(ctx context) error {
	o := s.domain.GetOper(ctx.path.Part(1))
	result, response := o.Run(ctx)

	if result != nil && result.Error() != nil {
		ctx.Debugf("Oper Failed: %+v", result)
		jsonResult, _ := json.Marshal(result)
		http.Error(ctx.res, string(jsonResult), http.StatusInternalServerError)
		return nil
	}

	ctx.Debugf("Oper Success")
	if response != nil {
		//encode response body
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return errors.Wrapf(err, "failed to encode response as JSON")
		}

		if strings.Contains(ctx.req.Header.Get("Accept"), "xml") {
			//convert JSON to XML
			responseObj, err := mxj.NewMapJson([]byte(jsonResponse))
			if err != nil {
				return errors.Wrapf(err, "failed to convert JSON response to object")
			}
			xmlBody, err := responseObj.Xml()
			if err != nil {
				return errors.Wrapf(err, "failed to encode response object as XML")
			}
			ctx.res.Header().Set("Content-Type", "application/xml")
			ctx.res.Write(xmlBody)
		} else {
			ctx.res.Header().Set("Content-Type", "application/json")
			ctx.res.Write(jsonResponse)
		}
	}
	return nil
} //server.serveOper()

type path struct {
	names []string
}

func pathSplit(p string) path {
	n := path{
		names: []string{},
	}
	pp := strings.Split(p, "/")
	for _, s := range pp {
		if len(s) > 0 {
			n.names = append(n.names, s)
		}
	}
	return n
}

func (p path) Part(n int) string {
	if n < 0 || n >= len(p.names) {
		return ""
	}
	return p.names[n]
}
