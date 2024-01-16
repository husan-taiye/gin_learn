package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
	//Send(ctx context.Context, id string, strings []string, num string) error
}
