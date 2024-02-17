package controller

import (
	"OpenTeens/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreateNoteHandler 创建笔记的Handler
func CreateNoteHandler(c *gin.Context) {
	var noteRequest struct {
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		AttachmentHash []string `json:"attachment_hash"`
		RealCreator    string   `json:"real_creator"`
		NoteTypeID     uint     `json:"note_type_id"`
		NeedCredit     bool     `json:"need_credit"`
		CostCredit     int      `json:"cost_credit"`
	}
	if err := c.BindJSON(&noteRequest); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}

	uploaderName := "" // 从上下文获取上传者姓名
	// 从上下文或其他地方获取上传者ID
	uploaderID := uint(0) // 示例ID，应该从上下文获取实际ID

	success, noteInfo := CreateNoteService(noteRequest.Title, noteRequest.Description, noteRequest.AttachmentHash, noteRequest.RealCreator, uploaderName, uploaderID, noteRequest.NoteTypeID, noteRequest.NeedCredit, noteRequest.CostCredit)

	if success {
		c.JSON(200, gin.H{"msg": "Note created successfully", "note_info": noteInfo})
	} else {
		c.JSON(500, gin.H{"msg": "Failed to create note"})
	}
}

// CreateNoteService 创建笔记的Service
func CreateNoteService(title, description string, attachmentHashes []string, realCreator, uploaderName string, uploaderID, noteTypeID uint, needCredit bool, costCredit int) (bool, *model.NoteInfo) {
	// 创建NoteInfo
	noteInfo := &model.NoteInfo{
		UploaderID:    uploaderID,
		RealCreator:   realCreator,
		NoteTitle:     title,
		NoteType:      strconv.Itoa(int(noteTypeID)),
		NoteStatus:    "等待审核",
		ReviewPass:    false,
		IfNeedCredit:  needCredit,
		CostCredit:    costCredit,
		NewestVersion: 1,
	}

	// 设置作者
	if realCreator == "" {
		noteInfo.RealCreator = uploaderName
	}

	// 创建NoteBody
	noteBody := model.NoteBody{
		NoteVersion:     1,
		NoteDescription: description,
	}
	for _, hash := range attachmentHashes {
		noteBody.Attachment = append(noteBody.Attachment, model.Attachment{FileHash: hash})
	}

	// 将NoteBody添加到NoteInfo
	noteInfo.NoteBodies = append(noteInfo.NoteBodies, noteBody)

	// 保存到数据库
	success := model.DBCreateNote(noteInfo)
	return success, noteInfo
}
