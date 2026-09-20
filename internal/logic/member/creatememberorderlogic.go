package member

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/member/pb"

	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMemberOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMemberOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMemberOrderLogic {
	return &CreateMemberOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMemberOrderLogic) CreateMemberOrder(req *types.CreateMemberOrderReq) (*types.CreateMemberOrderResp, error) {
	l.Infof("[CreateMemberOrder] level: %v, duration_days: %v, amount: %v", req.Level, req.DurationDays, req.Amount)
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	// 1. 调用 MemberRPC 创建订单 (在 rpc 中，金额被写死了，但这符合安全要求，我们只需传入 PlanId)
	// 如果有差价，后端目前由于写死价格无法处理，所以为了沙箱测试，这里先强制使用后端原价创建订单
	// 其实前端显示的是差价，后端支付宝支付的是按系统获取的价格
	
	createResp, err := l.svcCtx.MemberRPC.CreateOrder(l.ctx, &pb.CreateOrderRequest{
		UserId:       userId,
		Level:        req.Level,
		DurationDays: req.DurationDays,
		PayChannel:   req.PayChannel,
	})
	if err != nil {
		l.Errorf("[CreateMemberOrder] RPC err: %v", err)
		return nil, err
	}

	// 2. 生成支付宝支付链接
	payUrl := ""
	if l.svcCtx.AlipayClient != nil && req.PayChannel == "alipay" {
		var p = alipay.TradePagePay{}
		p.NotifyURL = l.svcCtx.Config.Alipay.NotifyURL
		p.ReturnURL = l.svcCtx.Config.Alipay.ReturnURL // 支付完成后的回调跳转地址
		p.Subject = "ThinkTalk 会员订阅"
		p.OutTradeNo = createResp.OrderSn
		p.TotalAmount = fmt.Sprintf("%.2f", float64(createResp.Amount)/100.0)
		p.ProductCode = "FAST_INSTANT_TRADE_PAY"

		url, err := l.svcCtx.AlipayClient.TradePagePay(p)
		if err != nil {
			l.Errorf("[CreateMemberOrder] Alipay generate URL err: %v", err)
			return nil, err
		}
		payUrl = url.String()
	}

	return &types.CreateMemberOrderResp{
		OrderSn: createResp.OrderSn,
		Amount:  createResp.Amount,
		PayUrl:  payUrl,
	}, nil
}
