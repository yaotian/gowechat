package pay

import (
	"github.com/yaotian/gowechat/mch/base"
	"github.com/yaotian/gowechat/server/context"
)

//Pay pay
type Pay struct {
	base.MchBase
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
