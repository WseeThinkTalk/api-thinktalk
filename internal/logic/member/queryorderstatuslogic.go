package member

import (
	"context"
	"errors"
	"fmt"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/pkg/code"
	"api-thinktalk/client/member/pb"

	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOrderStatusLogic {
	return &QueryOrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryOrderStatusLogic) QueryOrderStatus(req *types.QueryOrderStatusReq) (*types.QueryOrderStatusResp, error) {
	if req.OrderSn == "" {
		return nil, code.OrderSnEmpty
	}

	if l.svcCtx.AlipayClient == nil {
		l.Errorf("[QueryOrderStatus] AlipayClient is nil")
		return nil, errors.New("alipay client not initialized")
	}

	// 1. 主动查询支付宝交易状态
	tradeResp, err := l.svcCtx.AlipayClient.TradeQuery(l.ctx, alipay.TradeQuery{
		OutTradeNo: req.OrderSn,
	})
	if err != nil {
		l.Errorf("[QueryOrderStatus] TradeQuery err: %v, orderSn: %s", err, req.OrderSn)
		return nil, code.OrderQueryFailed
	}

	// 2. 构建返回信息
	resp := &types.QueryOrderStatusResp{
		OrderSn:    req.OrderSn,
		TradeNo:    tradeResp.TradeNo,
		TotalAmount: tradeResp.TotalAmount,
	}

	// 3. 根据支付宝交易状态处理
	switch tradeResp.TradeStatus {
	case alipay.TradeStatusSuccess, alipay.TradeStatusFinished:
		// 支付成功：主动调用 PayCallback 完成订单处理（幂等安全）
		l.Infof("[QueryOrderStatus] trade success from Alipay, calling PayCallback for orderSn: %s", req.OrderSn)
		_, err := l.svcCtx.MemberRPC.PayCallback(l.ctx, &pb.PayCallbackRequest{
			OrderSn:       req.OrderSn,
			TransactionId: tradeResp.TradeNo,
		})
		if err != nil {
			l.Errorf("[QueryOrderStatus] PayCallback RPC err: %v, orderSn: %s", err, req.OrderSn)
			return nil, fmt.Errorf("payment confirmed by Alipay but order processing failed: %w", err)
		}
		l.Infof("[QueryOrderStatus] Successfully processed payment for orderSn: %s", req.OrderSn)
		resp.Status = 1 // OrderStatusPaid
		resp.StatusText = "已支付"

	case alipay.TradeStatusWaitBuyerPay:
		resp.Status = 0 // OrderStatusPending
		resp.StatusText = "待支付"

	case alipay.TradeStatusClosed:
		resp.Status = 2 // OrderStatusRefunded / Closed
		resp.StatusText = "已关闭"

	default:
		// 未知的支付宝交易状态，返回待支付
		l.Infof("[QueryOrderStatus] unknown trade status: %s, orderSn: %s", tradeResp.TradeStatus, req.OrderSn)
		resp.Status = 0
		resp.StatusText = "待支付"
	}

	return resp, nil
}
