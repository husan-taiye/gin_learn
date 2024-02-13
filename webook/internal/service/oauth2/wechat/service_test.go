//go:build manual

package wechat

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestOAuth2Service_manual_VerifyCode(t *testing.T) {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_ID ")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET")
	}
	svc := NewService(appId, appKey)
	res, err := svc.VerifyCode(context.Background(), "051D6b000Yn4FQ14Rd300FgOF33D6b0s", "state")
	require.NoError(t, err)
	t.Log(res)
}
