// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"api-thinktalk/article/internal/svc"
	"api-thinktalk/article/internal/types"
	"api-thinktalk/client/article/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAuditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAuditLogic {
	return &AdminAuditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminAuditLogic) AdminAudit(req *types.AdminAuditRequest) (resp *types.AdminAuditResponse, err error) {
	_, err = l.svcCtx.ArticleRPC.AdminAudit(l.ctx, &article.AdminAuditRequest{
		ArticleId: req.ArticleId,
		Status:    req.Status,
	})
	if err != nil {
		l.Errorf("AdminAudit articleId: %d, status: %d error: %v", req.ArticleId, req.Status, err)
		return nil, err
	}

	return &types.AdminAuditResponse{}, nil
}
