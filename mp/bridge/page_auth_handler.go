package bridge

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/yaotian/gowechat/mp/oauth"
	"github.com/yaotian/gowechat/wxcontext"
)

//PageOAuthHandler 微信网页授权
type PageOAuthHandler struct {
	*wxcontext.Context
	oauth *oauth.Oauth

	oAuthURL     string
	urlNeedOAuth string

	openIDExisting          bool
	checkOpenIDExistingFunc func(openID string) bool

	oauth.UserInfo
	afterGetUserInfoFunc func(user oauth.UserInfo)
}

//NewPageOAuthHandler PageOAuthHandler初始化
func NewPageOAuthHandler(context *wxcontext.Context, oAuthURL string) *PageOAuthHandler {
	srv := new(PageOAuthHandler)
	srv.Context = context
	srv.oauth = new(oauth.Oauth)
	srv.oAuthURL = oAuthURL
	return srv
}

func (c *PageOAuthHandler) getCallbackURL() (u string) {
	return fmt.Sprintf("%s?target=%s", c.oAuthURL, url.QueryEscape(c.urlNeedOAuth))
}

//Handler handler
func (c *PageOAuthHandler) Handler() (err error) {
	code := c.Query("code")
	state := c.Query("state")
	c.urlNeedOAuth = c.Query("target")
	if code != "" {
		var acsTkn oauth.ResAccessToken
		acsTkn, err = c.oauth.GetUserAccessToken(code)
		if err != nil {
			return
		}
		openID := acsTkn.OpenID
		if c.checkOpenIDExistingFunc(openID) { //系统中已经存在openID
			http.Redirect(c.Writer, nil, c.urlNeedOAuth, 302)
			return
		}
		if state == "base" {
			c.oauth.Redirect(c.Writer, c.getCallbackURL(), "snsapi_userinfo", "userinfo")
			return
		}
		c.UserInfo, err = c.oauth.GetUserInfo(acsTkn.AccessToken, openID)
		if err == nil {
			c.afterGetUserInfoFunc(c.UserInfo)
			http.Redirect(c.Writer, nil, c.urlNeedOAuth, 302)
			return
		}
	}
	c.oauth.Redirect(c.Writer, c.getCallbackURL(), "snsapi_base", "base")
	return
}
