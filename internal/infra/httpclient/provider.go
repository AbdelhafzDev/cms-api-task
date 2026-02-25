package httpclient

import "go.uber.org/fx"


var Module = fx.Module("httpclient",
	fx.Provide(NewDefault),
)

func NewDefault() *Client {
	return New(nil)
}
