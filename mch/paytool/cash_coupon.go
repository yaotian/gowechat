package paytool

import (
	"github.com/uzicloud/gowechat/mch/base"
	"github.com/uzicloud/gowechat/wxcontext"
)

//PayTool pay tool
type PayTool struct {
	base.MchBase
}

//NewPayTool 实例化
func NewPayTool(context *wxcontext.Context) *PayTool {
	payT := new(PayTool)
	payT.Context = context
	return payT
}

//SendRedPackRaw 发现金红包
func (c *PayTool) SendRedPackRaw(req map[string]string) (resp map[string]string, err error) {
	return c.PostXML("https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack", req, true)
}
