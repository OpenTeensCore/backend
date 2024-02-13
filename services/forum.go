package services

import (
	"OpenTeens/model"
)

var (
	// listener dict {chatname: []listener}
	forumListeners = make(map[string][]func(string, model.ForumChatMsg))
)

// ForumCreateRoom 创建一个聊天室
func ForumCreateRoom(owner uint, chatname string) (bool, string) {
	ok := model.DBForumCreateChatRoom(owner, chatname)
	if ok {
		return true, chatname
	}
	return false, chatname
}

// ForumFetchHistory 获取历史聊天记录
func ForumFetchHistory(chatname string, start int, end int) (bool, []model.ForumChatMsg) {
	return model.DBForumGetChatMsgs(chatname, start, end)
}

// ForumWriteMsg 写入一条聊天消息
func ForumWriteMsg(chatname string, sender uint, content string) bool {
	// store
	ok, msg := model.DBForumCreateChatMsg(chatname, sender, content)
	if !ok {
		return false
	}

	// listener
	for _, listener := range forumListeners[chatname] {
		if listener != nil {
			listener(chatname, msg)
		}
	}

	return ok
}

// ForumAddListener 添加一个聊天室监听器，返回监听器序号
func ForumAddListener(chatname string, listener func(string, model.ForumChatMsg)) uint {
	forumListeners[chatname] = append(forumListeners[chatname], listener)

	return uint(len(forumListeners[chatname]) - 1)
}

// ForumRemoveListener 移除一个聊天室监听器
func ForumRemoveListener(chatname string, listenerID uint) bool {
	forumListeners[chatname][listenerID] = nil

	return true
}
