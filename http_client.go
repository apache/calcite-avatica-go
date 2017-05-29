package avatica

import (
	"bytes"
	"io/ioutil"
	"net/http"

	avaticaMessage "github.com/Boostport/avatica/message"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/xinsnake/go-http-digest-auth-client"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type httpClientAuthConfig struct {
	username           string
	password           string
	authenticationType authentication
}

// httpClient wraps the default http.Client to communicate with the Avatica server.
type httpClient struct {
	host       string
	authConfig httpClientAuthConfig

	httpClient *http.Client
}

// NewHTTPClient creates a new httpClient from a host.
func NewHTTPClient(host string, authenticationConf httpClientAuthConfig) *httpClient {

	client := cleanhttp.DefaultPooledClient()

	if authenticationConf.authenticationType == digest {
		rt := digest_auth_client.NewTransport(authenticationConf.username, authenticationConf.password)
		client.Transport = &rt
	}

	return &httpClient{
		host:       host,
		authConfig: authenticationConf,

		httpClient: client,
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

	req, err := http.NewRequest("POST", c.host, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-google-protobuf")

	if c.authConfig.authenticationType == basic {
		req.SetBasicAuth(c.authConfig.username, c.authConfig.password)
	}

	res, err := ctxhttp.Do(ctx, c.httpClient, req)

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
