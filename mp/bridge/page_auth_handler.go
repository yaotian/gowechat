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
	// *wxcontext.Context
	*oauth.Oauth

	oAuthCallbackURL string
	urlNeedOAuth     string

	openID                  string
	openIDExisting          bool
	checkOpenIDExistingFunc func(openID string) bool

	oauth.UserInfo
	afterGetUserInfoFunc func(user oauth.UserInfo)
}

//NewPageOAuthHandler PageOAuthHandler初始化
func NewPageOAuthHandler(context *wxcontext.Context, oAuthCallbackURL string) *PageOAuthHandler {
	pa := new(PageOAuthHandler)
	pa.Oauth = oauth.NewOauth(context)
	pa.oAuthCallbackURL = oAuthCallbackURL
	return pa
}

func (c *PageOAuthHandler) getCallbackURL() (u string) {
	return fmt.Sprintf("%s?target=%s", c.oAuthCallbackURL, url.QueryEscape(c.urlNeedOAuth))
}

//SetFuncCheckOpenIDExisting 设置检查OpenID在您的系统中是否已经存在
func (c *PageOAuthHandler) SetFuncCheckOpenIDExisting(handler func(string) bool) {
	c.checkOpenIDExistingFunc = handler
}

//SetFuncAfterGetUserInfo 设置获得用户信息后执行
func (c *PageOAuthHandler) SetFuncAfterGetUserInfo(handler func(oauth.UserInfo)) {
	c.afterGetUserInfoFunc = handler
}

//Handle handler
func (c *PageOAuthHandler) Handle() (err error) {
	code := c.Query("code")
	state := c.Query("state")
	c.urlNeedOAuth = c.Query("target")
	if code != "" {
		var acsTkn oauth.ResAccessToken
		acsTkn, err = c.GetUserAccessToken(code)
		if err != nil {
			return
		}
		openID := acsTkn.OpenID
		if c.checkOpenIDExistingFunc(openID) { //系统中已经存在openID
			http.Redirect(c.Writer, c.Request, c.urlNeedOAuth, 302)
			return
		}
		if state == "base" {
			c.Redirect(c.getCallbackURL(), "snsapi_userinfo", "userinfo")
			return
		}
		c.UserInfo, err = c.GetUserInfo(acsTkn.AccessToken, openID)
		if err == nil {
			c.afterGetUserInfoFunc(c.UserInfo)
			http.Redirect(c.Writer, c.Request, c.urlNeedOAuth, 302)
			return
		}
	}
	c.Redirect(c.getCallbackURL(), "snsapi_base", "base")
	return
}
