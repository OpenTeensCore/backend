package model

import (
	"OpenTeens/dao"
	"time"

	"gorm.io/gorm"
)

// ForumChatMsg 论坛聊天消息
type ForumChatMsg struct {
	gorm.Model
	Sender  uint
	Content string
}

// ForumChatRoom 论坛聊天
type ForumChatRoom struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	RoomURI  string `gorm:"primarykey"`
	Owner    uint
	Member   []uint
	Messages []uint // store message id
}

// DBForumCreateChatRoom 在数据库中创建一个聊天室
func DBForumCreateChatRoom(ownerID uint, roomURI string) bool {
	// Create a table
	chat := ForumChatRoom{
		RoomURI: roomURI,
		Owner:   ownerID,
		Member:  []uint{ownerID},
	}
	result := dao.DB.Create(&chat)

	return result.Error == nil
}

// DBForumCreateChatMsg 在数据库中创建一个聊天消息（消息存储）
func DBForumCreateChatMsg(roomURI string, senderID uint, content string) (bool, ForumChatMsg) {
	// Create a table
	msg := ForumChatMsg{
		Sender:  senderID,
		Content: content,
	}
	result := dao.DB.Create(&msg)

	// Update the chat room
	var chat ForumChatRoom
	dao.DB.Where("room_uri = ?", roomURI).First(&chat)
	chat.Messages = append(chat.Messages, msg.ID)

	return result.Error == nil, msg
}

// DBForumGetChatRoom 获取聊天室信息
func DBForumGetChatRoom(roomURI string) (bool, ForumChatRoom) {
	var chat ForumChatRoom
	result := dao.DB.Where("room_uri = ?", roomURI).First(&chat)

	return result.Error == nil, chat
}

// DBForumGetChatMsg 获取聊天消息
func DBForumGetChatMsg(msgID uint) (bool, ForumChatMsg) {
	var msg ForumChatMsg
	result := dao.DB.Where("id = ?", msgID).First(&msg)

	return result.Error == nil, msg
}

// DBForumGetChatMsgs （历史消息拉取）获取某一聊天室 start~end 的消息，为负数则表示从最新的消息开始
func DBForumGetChatMsgs(roomURI string, start int, end int) (bool, []ForumChatMsg) {
	var chat ForumChatRoom
	result := dao.DB.Where("room_uri = ?", roomURI).First(&chat)

	if result.Error != nil {
		return false, nil
	}

	var msgs []ForumChatMsg
	if start < 0 {
		start = len(chat.Messages) + start
	}
	if end < 0 {
		end = len(chat.Messages) + end
	}

	for i := start; i <= end; i++ {
		ok, msg := DBForumGetChatMsg(chat.Messages[i])
		if !ok {
			return false, nil
		}
		msgs = append(msgs, msg)
	}

	return true, msgs
}
