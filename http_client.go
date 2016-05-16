package avatica

import (
	"bytes"
	avaticaMessage "github.com/Boostport/avatica/message"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
)

// httpClient wraps the default http.Client to communicate with the Avatica server.
type httpClient struct {
	host       string
	httpClient *http.Client
}

// NewHTTPClient creates a new httpClient from a host.
func NewHTTPClient(host string) *httpClient {

	return &httpClient{
		host:       host,
		httpClient: cleanhttp.DefaultPooledClient(),
	}
}

// post posts a protocol buffer message to the Avatica server.
func (c *httpClient) post(ctx context.Context, message proto.Message) (proto.Message, error) {

	wrapped, err := proto.Marshal(message)

	if err != nil {
		return nil, err
	}

	wire := &avaticaMessage.WireMessage{
		Name:           classNameFromRequest(message),
		WrappedMessage: wrapped,
	}

	body, err := proto.Marshal(wire)

	if err != nil {
		return nil, err
	}

	res, err := ctxhttp.Post(ctx, c.httpClient, c.host, "application/x-google-protobuf", bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	result := &avaticaMessage.WireMessage{}

	err = proto.Unmarshal(response, result)

	if err != nil {
		return nil, err
	}

	inner, err := responseFromClassName(result.Name)

	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(result.WrappedMessage, inner)

	if err != nil {
		return nil, err
	}

	if v, ok := inner.(*avaticaMessage.ErrorResponse); ok {
		return nil, errorResponseToResponseError(v)
	}

	return inner, nil
}
