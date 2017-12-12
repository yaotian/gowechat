package pay

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/yaotian/gowechat/server/context"
)

//Pay pay
type Pay struct {
	*context.Context
}

//NewPay 实例化
func NewPay(context *context.Context) *Pay {
	pay := new(Pay)
	pay.Context = context
	return pay
}

// UnifiedOrder 统一下单.
func (c *Pay) UnifiedOrder(req map[string]string) (resp map[string]string, err error) {
	return c.PostXML("https://api.mch.weixin.qq.com/pay/unifiedorder", req, false)
}

// 查询订单.
func (c *Pay) OrderQuery(req map[string]string) (resp map[string]string, err error) {
	return c.PostXML("https://api.mch.weixin.qq.com/pay/orderquery", req, false)
}

// 关闭订单.
func (c *Pay) CloseOrder(req map[string]string) (resp map[string]string, err error) {
	return c.PostXML("https://api.mch.weixin.qq.com/pay/closeorder", req, false)
}

// 申请退款.
//  NOTE: 请求需要双向证书.
func (c *Pay) Refund(req map[string]string) (resp map[string]string, err error) {
	return c.PostXML("https://api.mch.weixin.qq.com/secapi/pay/refund", req, true)
}

// 查询退款.
func (c *Pay) RefundQuery(req map[string]string) (resp map[string]string, err error) {
	return c.PostXML("https://api.mch.weixin.qq.com/pay/refundquery", req, false)
}

//PostXML postXML
func (c *Pay) PostXML(url string, req map[string]string, needSSL bool) (resp map[string]string, err error) {
	bodyBuf := textBufferPool.Get().(*bytes.Buffer)
	bodyBuf.Reset()
	defer textBufferPool.Put(bodyBuf)

	if err = FormatMapToXML(bodyBuf, req); err != nil {
		return
	}

	//需要ssl，就需要ssl client
	client := c.HTTPClient
	if needSSL {
		client = c.SHTTPClient
	}

	httpResp, err := client.Post(url, "text/xml; charset=utf-8", bodyBuf)
	if err != nil {
		return resp, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}

	if resp, err = ParseXMLToMap(bytes.NewReader(respBody)); err != nil {
		return
	}

	beego.Debug(resp)

	// 判断协议状态
	ReturnCode, ok := resp["return_code"]
	if !ok {
		err = errors.New("no return_code parameter")
		return
	}
	if ReturnCode != ReturnCodeSuccess {
		err = &Error{
			ReturnCode: ReturnCode,
			ReturnMsg:  resp["return_msg"],
		}
		return
	}

	// 安全考虑, 做下验证
	mchId, ok := resp["mch_id"]
	if ok && mchId != c.MchID {
		err = fmt.Errorf("mch_id mismatch, have: %q, want: %q", mchId, c.MchID)
		return
	}

	//发送红包的情况，不需要验证这些，因为有的信息没有
	if !needSSL {
		appId, ok := resp["appid"]
		if ok && appId != c.AppID {
			err = fmt.Errorf("appid mismatch, have: %q, want: %q", appId, c.AppID)
			return
		}

		// 认证签名
		signature1, ok := resp["sign"]
		if !ok {
			err = errors.New("no sign parameter")
			return
		}
		signature2 := Sign(resp, c.MchAPIKey, nil)
		if signature1 != signature2 {
			err = fmt.Errorf("check signature failed, \r\ninput: %q, \r\nlocal: %q", signature1, signature2)
			return
		}
	}
	return
}
