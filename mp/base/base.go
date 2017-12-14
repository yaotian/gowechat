package base

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yaotian/gowechat/context"
	"github.com/yaotian/gowechat/util"
)

//MpBase 微信公众平台,基本类
type MpBase struct {
	*context.Context
}

//HTTPGetWithAccessToken 微信公众平台中，自动加上access_token变量的GET调用，
//如果失败，会清空AccessToken cache, 再试一次
func (c *MpBase) HTTPGetWithAccessToken(url string) (resp []byte, err error) {
	retry := 1
Do:
	var accessToken string
	accessToken, err = c.GetAccessToken()
	if err != nil {
		return
	}

	var target = ""
	if strings.Contains(url, "?") {
		target = fmt.Sprintf("%s&access_token=%s", url, accessToken)
	} else {
		target = fmt.Sprintf("%s?access_token=%s", url, accessToken)
	}

	var reponse *http.Response
	reponse, err = http.Get(target)
	if err != nil {
		return
	}
	defer reponse.Body.Close()

	resp, err = ioutil.ReadAll(reponse.Body)
	err = util.CheckCommonError(resp)
	if err == util.ErrUnmarshall {
		return
	}
	if err != nil {
		if retry > 0 {
			retry--
			c.CleanAccessTokenCache()
			goto Do
		}
		return
	}
	return
}
