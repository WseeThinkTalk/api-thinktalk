package code

import "api-thinktalk/pkg/xcode"

var (
	// User & Auth
	RegisterNameEmpty     = xcode.New(100004, "名字不能为空")
	RegisterMobileEmpty   = xcode.New(10001, "注册手机号不能为空")
	VerificationCodeEmpty = xcode.New(100002, "验证码不能为空")
	MobileHasRegistered   = xcode.New(100003, "手机号已经注册")
	LoginMobileEmpty      = xcode.New(100003, "手机号不能为空")
	RegisterPasswordEmpty = xcode.New(100004, "密码不能为空")

	// Article
	ParseFormErr              = xcode.New(30000, "解析表单失败")
	GetBucketErr              = xcode.New(30001, "获取bucket实例失败")
	PutBucketErr              = xcode.New(30002, "上传bucket失败")
	GetObjDetailErr           = xcode.New(30003, "获取对象详细信息失败")
	ArtitleTitleEmpty         = xcode.New(30004, "文章标题为空")
	ArticleContentTooFewWords = xcode.New(30005, "文章内容字数太少")
	ArticleCoverEmpty         = xcode.New(30006, "文章封面为空")

	// Member & Order
	UserIdEmpty           = xcode.New(90001, "用户ID不能为空")
	LevelInvalid          = xcode.New(90002, "会员等级无效")
	TransactionIdEmpty    = xcode.New(90003, "支付流水号不能为空")
	OrderNotFound         = xcode.New(90004, "订单不存在")
	MemberNotFound        = xcode.New(90005, "会员信息不存在")
	DuplicateTransaction  = xcode.New(90006, "重复的支付流水号")
	OrderSnEmpty          = xcode.New(90007, "订单号不能为空")
	OrderAlreadyProcessed = xcode.New(90008, "订单已被处理，无法重复更新")
	OrderQueryFailed      = xcode.New(90009, "主动查询订单状态失败")
)

