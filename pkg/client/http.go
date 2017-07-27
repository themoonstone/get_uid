package client

import (
	"context"
	//	"encoding/json"

	"io"
	"log"
	"net/http"
	"net/url"

	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	path = "/user/getuid"
)

type ReqClient struct {
	url *url.URL
}

type Response struct {
	Body   io.ReadCloser
	String string
}

func NewHttpClient(host string) *ReqClient {

	return &ReqClient{
		&url.URL{
			Scheme: "http",
			Host:   host,
		},
	}
}

func (c *ReqClient) SendRequest() (interface{}, error) {
	var (
		encode = func(context.Context, *http.Request, interface{}) error { return nil }

		decode = DecodeUidResponseFunc
	)
	c.url.Path = path
	client := httptransport.NewClient(
		http.MethodGet,
		c.url,
		encode,
		decode,
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

/*
type Uid struct {
	id         int64  `json:"id"`
	machine_id uint16 `json:"machine-id"`
	msb        int64  `json:"msb"`
	sequence   uint16 `json:"sequence"`
	time       uint64 `json:"time"`
}
*/
func DecodeUidResponseFunc(_ context.Context, r *http.Response) (interface{}, error) {
	buffer := make([]byte, 64)
	r.Body.Read(buffer)
	return Response{r.Body, string(buffer)}, nil
}
