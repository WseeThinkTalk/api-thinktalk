package member

import (
	"context"
	"strconv"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	member "api-thinktalk/client/member/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberOrderListLogic {
	return &MemberOrderListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MemberOrderListLogic) MemberOrderList(userId int64, req *types.MemberOrderListRequest) (resp *types.MemberOrderListResponse, err error) {
	resp = new(types.MemberOrderListResponse)

	rpcResp, err := l.svcCtx.MemberRPC.MemberOrderList(l.ctx, &member.MemberOrderListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[MemberOrderList] rpc err: %v", err)
		return nil, err
	}

	var items []*types.MemberOrderItem
	if rpcResp != nil && rpcResp.Data != nil {
		for _, v := range rpcResp.Data.Items {
			if v != nil {
				items = append(items, &types.MemberOrderItem{
					OrderId:        strconv.FormatInt(v.Id, 10),
					UserId:         v.UserId,
					TargetLevel:    v.Level,
					DurationMonths: v.DurationDays / 30,
					Amount:         float64(v.Amount) / 100.0,
					PayChannel:     v.PayChannel,
					PayStatus:      v.Status,
					CreateTime:     v.CreateTime,
				})
			}
		}
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	resp.Items = items
	return resp, nil
}
