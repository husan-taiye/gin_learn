package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_learn/webook/internal/domain"
	"net/http"
	"net/url"
)

var redirectURI = "https://xxx.com/oauth2/wechat/callback"

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type OAuth2Service struct {
	appId     string
	appSecret string
	client    *http.Client
}

// NewServiceV1 不偷懒
func NewServiceV1(appId string, appSecret string, client *http.Client) Service {
	return &OAuth2Service{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
	}
}

func NewService(appId string, appSecret string) Service {
	return &OAuth2Service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (o *OAuth2Service) AuthURL(ctx context.Context, state string) (string, error) {
	const wxAuthUrlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	url.PathEscape(redirectURI)
	return fmt.Sprintf(wxAuthUrlPattern, o.appId, redirectURI, state), nil
}

func (o *OAuth2Service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const verifyCodeUrlPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(verifyCodeUrlPattern, o.appId, o.appSecret, code)
	//resp, err := http.Get(target)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	//req, err := http.NewRequest(http.MethodGet, target, nil)
	//// 会产生复制, 性能差
	//req := req.WithContext(ctx)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)

	// 整个响应都读出来 不推荐，因为unmarshal会再读一遍 跟readall合计两遍
	//body, err := io.ReadAll(resp.Body)
	//err = json.Unmarshal(body, &res)

	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("微信返回错误响应， 错误码：%d，错误信息：%s ", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionID,
	}, nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AssessToken  string `json:"assess_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenId  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}
