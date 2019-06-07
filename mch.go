package gowechat

import (
	"github.com/uzicloud/gowechat/mch/pay"
	"github.com/uzicloud/gowechat/mch/paytool"
)

//MchMgr mch mgt
type MchMgr struct {
	*Wechat
}

// GetPay 基本支付api
func (wc *MchMgr) GetPay() *pay.Pay {
	return pay.NewPay(wc.Context)
}

// GetPayTool 支付工具，发红包等
func (wc *MchMgr) GetPayTool() *paytool.PayTool {
	return paytool.NewPayTool(wc.Context)
}
