package member

import (
	"context"

	"api-thinktalk/client/member/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

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

	rpcResp, err := l.svcCtx.MemberRPC.MemberOrderList(l.ctx, &pb.MemberOrderListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[MemberOrderList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.MemberOrderItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		items = append(items, &types.MemberOrderItem{
			Id:           item.Id,
			UserId:       item.UserId,
			Level:        item.Level,
			DurationDays: item.DurationDays,
			Amount:       item.Amount,
			PayChannel:   item.PayChannel,
			Status:       item.Status,
			CreateTime:   item.CreateTime,
		})
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
