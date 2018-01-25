package pay

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/yaotian/gowechat/mch/base"
	"github.com/yaotian/gowechat/util"
)

//官方文档： https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=13_4&index=3

const (
	//SceneIDPromotion 商品促销
	SceneIDPromotion = "PRODUCT_1"

	//SceneIDLuckyDraw 抽奖
	SceneIDLuckyDraw = "PRODUCT_2"

	//SceneIDPrize 虚拟物品兑奖
	SceneIDPrize = "PRODUCT_3"

	//SceneIDBenefit 企业内部福利
	SceneIDBenefit = "PRODUCT_4"

	//SceneIDAgentBonous 渠道分润
	SceneIDAgentBonous = "PRODUCT_5"

	//SceneIDAgentInsurance 保险回馈
	SceneIDAgentInsurance = "PRODUCT_6"

	//SceneIDLottery 彩票派奖
	SceneIDLottery = "PRODUCT_7"

	//SceneIDTax 税务刮奖
	SceneIDTax = "PRODUCT_8"
)

//RedPackInput 发红包的配置
type RedPackInput struct {
	ToOpenID string //接红包的OpenID
	MoneyFen int    //分为单位

	SendName string //商户名称，String(32) 谁发的红包，一般为发红包的单位
	Wishing  string //红包祝福语 String(128) “感谢您参加猜灯谜活动，祝您元宵节快乐！”
	ActName  string //活动名称 String(32) 猜灯谜抢红包活动
	Remark   string //备注 String(256)

	IP string

	//非必填，但大于200元，此必填, 有8个选项可供选择
	SceneID string
}

func (m *RedPackInput) Check() (isGood bool, err error) {
	if input.ToOpenID == "" || input.MoneyFen == 0 || input.SendName == "" || input.Wishing == "" || input.ActName == "" || input.Remark == "" || input.IP == "" {
		err = fmt.Errorf("%s", "Input有必填项没有值")
		return
	}

	if input.MoneyFen >= 200*100 && input.SceneID == "" {
		err = fmt.Errorf("%s", "大于200元的红包，必须设置SceneID")
		return
	}
	return true, nil
}

//SendRedPack 发红包
func (c *PayTool) SendRedPack(input RedPackInput) (isSuccess bool, err error) {
	if isGood, err := input.Check(); !isGood {
		return false, err
	}

	now := time.Now()
	dayStr := beego.Date(now, "Ymd")

	billno := c.MchID + dayStr + util.RandomStr(10)

	var signMap = make(map[string]string)
	signMap["nonce_str"] = util.RandomStr(5)
	signMap["mch_billno"] = billno //mch_id+yyyymmdd+10位一天内不能重复的数字
	signMap["mch_id"] = c.MchID
	signMap["wxappid"] = c.AppID
	signMap["send_name"] = input.SendName
	signMap["re_openid"] = input.ToOpenID
	signMap["total_amount"] = util.ToStr(input.MoneyFen)
	signMap["total_num"] = "1"
	signMap["wishing"] = input.Wishing
	signMap["client_ip"] = input.IP
	signMap["act_name"] = input.ActName
	signMap["remark"] = input.Remark
	signMap["sign"] = base.Sign(signMap, c.MchAPIKey, nil)

	respMap, err := c.SendRedPack(signMap)
	if err != nil {
		return false, err
	}

	result_code, ok := respMap["result_code"]
	if !ok {
		err = errors.New("no result_code")
		return false, err
	}

	if result_code != "SUCCESS" {
		err = errors.New("result code is not success")
		return false, err
	}

	mch_billno, ok := respMap["mch_billno"]
	if !ok {
		err = errors.New("no mch_billno")
		return false, err
	}

	if billno != mch_billno {
		err = errors.New("billno is not correct")
		beego.Error(err)
		return false, err
	}

	return true, nil
}
