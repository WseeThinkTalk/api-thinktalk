package member

import (
	"context"
	"errors"
	"net/http"

	"api-thinktalk/internal/svc"
	member "api-thinktalk/client/member/pb"

	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type AlipayNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlipayNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlipayNotifyLogic {
	return &AlipayNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlipayNotifyLogic) AlipayNotify(req *http.Request) error {
	if l.svcCtx.AlipayClient == nil {
		l.Errorf("[AlipayNotify] AlipayClient is nil")
		return errors.New("alipay client not initialized")
	}

	// 1. 解析支付宝异步通知请求并验签
	noti, err := l.svcCtx.AlipayClient.GetTradeNotification(req)
	if err != nil {
		l.Errorf("[AlipayNotify] GetTradeNotification err: %v", err)
		return err
	}

	// 2. 只有交易成功才处理
	if noti.TradeStatus == alipay.TradeStatusSuccess {
		// 3. 调用 MemberRPC 的 PayCallback 完成订单状态更新与会员权益下发
		_, err := l.svcCtx.MemberRPC.PayCallback(l.ctx, &member.PayCallbackRequest{
			OrderSn:       noti.OutTradeNo,
			TransactionId: noti.TradeNo,
		})
		if err != nil {
			l.Errorf("[AlipayNotify] PayCallback RPC err: %v, orderSn: %s", err, noti.OutTradeNo)
			return err
		}
		l.Infof("[AlipayNotify] Successfully processed payment for orderSn: %s", noti.OutTradeNo)
	}

	return nil
}
