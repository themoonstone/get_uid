package client // import "get_uid/pkg/client"

func DecodeUidResponseFunc(_ context.Context, r *http.Response) (interface{}, error)
type ReqClient struct{ ... }
    func NewHttpClient(host string) *ReqClient
type Response struct{ ... }
