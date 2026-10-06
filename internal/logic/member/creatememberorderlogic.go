package member

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	member "api-thinktalk/client/member/pb"

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

func (l *CreateMemberOrderLogic) CreateMemberOrder(req *types.CreateMemberOrderReq) (resp *types.CreateMemberOrderResp, err error) {
	resp = new(types.CreateMemberOrderResp)

	l.Infof("[CreateMemberOrder] level: %v, duration_days: %v, amount: %v", req.Level, req.DurationDays, req.Amount)
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	createResp, err := l.svcCtx.MemberRPC.CreateOrder(l.ctx, &member.CreateOrderRequest{
		UserId:       userId,
		Level:        req.Level,
		DurationDays: req.DurationDays,
		Amount:       req.Amount,
		PayChannel:   req.PayChannel,
	})
	if err != nil {
		l.Errorf("[CreateMemberOrder] RPC err: %v", err)
		return nil, err
	}

	orderSn := ""
	var orderAmount int64
	if createResp != nil && createResp.Data != nil {
		orderSn = createResp.Data.OrderSn
		orderAmount = createResp.Data.Amount
	}

	// 2. 生成支付宝支付链接
	payUrl := ""
	if l.svcCtx.AlipayClient != nil && req.PayChannel == "alipay" {
		var p = alipay.TradePagePay{}
		p.NotifyURL = l.svcCtx.Config.Alipay.NotifyURL
		p.ReturnURL = l.svcCtx.Config.Alipay.ReturnURL
		p.Subject = "ThinkTalk 会员订阅"
		p.OutTradeNo = orderSn
		p.TotalAmount = fmt.Sprintf("%.2f", float64(orderAmount)/100.0)
		p.ProductCode = "FAST_INSTANT_TRADE_PAY"

		url, err := l.svcCtx.AlipayClient.TradePagePay(p)
		if err != nil {
			l.Errorf("[CreateMemberOrder] Alipay generate URL err: %v", err)
			return nil, err
		}
		payUrl = url.String()
	}

	resp.OrderSn = orderSn
	resp.Amount = orderAmount
	resp.PayUrl = payUrl
	return resp, nil
}
