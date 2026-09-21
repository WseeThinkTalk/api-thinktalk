// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"api-thinktalk/pkg/code"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/pkg/encrypt"
	"api-thinktalk/pkg/jwt"
	"api-thinktalk/pkg/xcode"
	"context"
	"strings"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	resp = new(types.LoginResponse)

	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.LoginMobileEmpty
	}
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	req.Password = strings.TrimSpace(req.Password)

	var isPassMode bool
	if len(req.VerificationCode) > 0 {
		isPassMode = false
	} else if len(req.Password) > 0 {
		isPassMode = true
	} else {
		return nil, code.VerificationCodeEmpty
	}

	if !isPassMode {
		if err := CheckVerificationCode(l.svcCtx.RDB, req.Mobile, req.VerificationCode); err != nil {
			logx.Errorf("checkVerificationCode error: %v", err)
			return nil, err
		}
	}

	encMobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("encMobile error: %v", err)
		return nil, err
	}
	mobile, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{Mobile: encMobile})
	if err != nil {
		logx.Errorf("findByMobile error: %v", err)
		return nil, err
	}
	if mobile == nil || mobile.UserId == 0 {
		return nil, xcode.AccessDenied
	}

	if isPassMode {
		valid, needsUpgrade := encrypt.VerifyPassword(req.Password, mobile.Password)
		if !valid {
			return nil, xcode.AccessDenied
		}
		if needsUpgrade {
			newHash, hashErr := encrypt.HashPassword(req.Password)
			if hashErr == nil {
				_, _ = l.svcCtx.UserRPC.UpgradePassword(l.ctx, &user.UpgradePasswordRequest{
					UserId:       mobile.UserId,
					PasswordHash: newHash,
				})
			}
		}
	}

	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": mobile.UserId,
		},
	})
	if err != nil {
		logx.Errorf("buildTokens error: %v", err)
		return nil, err
	}

	if !isPassMode {
		_ = deleteActivationCode(req.Mobile, req.VerificationCode, l.svcCtx.RDB)
	}

	resp.UserId = mobile.UserId
	resp.Token = types.Token{
		AccessToken:  token.AccessToken,
		AccessExpire: token.AccessExpire,
	}
	return resp, nil
}
