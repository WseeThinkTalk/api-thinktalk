// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	user "api-thinktalk/client/user/service"
	"api-thinktalk/pkg/util"
	"context"
	"fmt"
	"strconv"
	"time"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
	prefixActivation        = "biz#activation#%s"
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	resp = new(types.VerificationResponse)

	count, err := l.GetVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	if count > verificationLimitPerDay {
		logx.Errorf("mobile: %s verification count over limit", req.Mobile)
		return nil, err
	}
	code, err := GetActivationCode(req.Mobile, l.svcCtx.RDB)
	if err != nil {
		logx.Errorf("getActivationCode mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	if len(code) == 0 {
		code = util.GenerateCode(6)
	}
	
	// 为了本地测试/沙箱环境，这里直接把验证码打印到日志中，方便在控制台查看并填入
	logx.Infof("==================================================")
	logx.Infof("【系统提示】已为您生成手机号 %s 的验证码: %s", req.Mobile, code)
	logx.Infof("==================================================")

	_, err = l.svcCtx.UserRPC.SendSms(l.ctx, &user.SendSmsRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		logx.Errorf("sendSms mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	err = SaveActivationCode(req.Mobile, code, l.svcCtx.RDB)
	if err != nil {
		logx.Errorf("saveActivationCode mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	err = l.IncrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	return resp, nil
}

func (l *VerificationLogic) GetVerificationCount(moblie string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, moblie)
	get, err := l.svcCtx.RDB.Get(key)
	if err != nil {
		return 0, err
	}
	if len(get) == 0 {
		return 0, nil
	}
	return strconv.Atoi(get)
}

func (l *VerificationLogic) IncrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.RDB.Incr(key)
	if err != nil {
		return err
	}
	expireTime := int(util.EndOfDay(time.Now()).Unix())
	return l.svcCtx.RDB.Expire(key, expireTime)
}

func GetActivationCode(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, mobile)
	val, err := rds.Get(key)
	if err != nil {
		return "", nil
	}
	return val, nil
}
func SaveActivationCode(mobile string, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}

func deleteActivationCode(mobile string, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	_, err := rds.Del(key)
	return err
}
