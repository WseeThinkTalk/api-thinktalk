package social

import (
	reply "api-thinktalk/client/reply/pb"
	"api-thinktalk/internal/types"
)

func convertReplyItem(pb *reply.ReplyItem, getUserInfo func(int64) (string, string)) *types.ReplyItem {
	if pb == nil {
		return nil
	}
	name, avatar := getUserInfo(pb.ReplyUserId)
	item := &types.ReplyItem{
		Id:         pb.ReplyId,
		BizId:      pb.BizId,
		ObjId:      pb.TargetId,
		UserId:     pb.ReplyUserId,
		Username:   name,
		Avatar:     avatar,
		Content:    pb.Content,
		ParentId:   pb.ParentId,
		LikeNum:    pb.LikeNum,
		CreateTime: pb.CreateTime,
	}
	// 递归转换子评论列表
	for _, v := range pb.SubReplies {
		item.SubReplies = append(item.SubReplies, convertReplyItem(v, getUserInfo))
	}
	return item
}
