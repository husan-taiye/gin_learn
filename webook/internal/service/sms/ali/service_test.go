package ali

import (
	"context"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"testing"
)

func TestService_Send(t *testing.T) {
	type fields struct {
		templateCode *string
		signName     *string
		client       *dysmsapi.Client
	}
	type args struct {
		ctx            context.Context
		phoneNumber    string
		signName       string
		templateCode   string
		templateParams string
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
			s := Service{
				templateCode: tt.fields.templateCode,
				signName:     tt.fields.signName,
				client:       tt.fields.client,
			}
			if err := s.Send(tt.args.ctx, tt.args.phoneNumber, tt.args.signName, tt.args.templateCode, tt.args.templateParams); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
