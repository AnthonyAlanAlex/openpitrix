package common

import (
	"log"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "openpitrix.io/openpitrix/test/client"
)

type IgnoreLogger struct{}

func (IgnoreLogger) Printf(format string, args ...interface{}) {
}

func (IgnoreLogger) Debugf(format string, args ...interface{}) {
}

func GetClient() *apiclient.Openpitrix {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	transport := httptransport.New("localhost:9100", "/", nil)
	transport.SetDebug(true)
	transport.SetLogger(IgnoreLogger{})
	Client := apiclient.New(transport, strfmt.Default)
	return Client
}
