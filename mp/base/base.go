package base

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yaotian/gowechat/server/context"
	"github.com/yaotian/gowechat/util"
)

//MpBase 微信公众号
type MpBase struct {
	*context.Context
}

//HTTPGetWithAccessToken http get
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
