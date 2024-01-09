package tencent

import (
	"context"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"testing"
)

func TestService_Send(t *testing.T) {
	type fields struct {
		appId    *string
		SignName *string
		client   *sms.Client
	}
	type args struct {
		ctx     context.Context
		tpl     string
		args    []string
		numbers []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				appId:    tt.fields.appId,
				SignName: tt.fields.SignName,
				client:   tt.fields.client,
			}
			if err := s.Send(tt.args.ctx, tt.args.tpl, tt.args.args, tt.args.numbers...); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
