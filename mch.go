package gowechat

import (
	"github.com/yaotian/gowechat/mch/pay"
	"github.com/yaotian/gowechat/mch/paytool"
)

// GetPay 基本支付api
func (wc *Wechat) GetPay() *pay.Pay {
	return pay.NewPay(wc.Context)
}

// GetPayTool 支付工具，发红包等
func (wc *Wechat) GetPayTool() *paytool.PayTool {
	return paytool.NewPayTool(wc.Context)
}
