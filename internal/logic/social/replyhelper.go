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
		ReplyId:         pb.ReplyId,
		BizId:           pb.BizId,
		TargetId:        pb.TargetId,
		ReplyUserId:     pb.ReplyUserId,
		BeReplyUserId:   pb.BeReplyUserId,
		ParentId:        pb.ParentId,
		Content:         pb.Content,
		LikeNum:         pb.LikeNum,
		CreateTime:      pb.CreateTime,
		ReplyUserName:   name,
		ReplyUserAvatar: avatar,
	}
	for _, sub := range pb.SubReplies {
		item.SubReplies = append(item.SubReplies, convertReplyItem(sub, getUserInfo))
	}
	return item
}
