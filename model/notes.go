package model

import (
	"OpenTeens/dao"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type NoteInfo struct {
	gorm.Model
	UploaderID    uint
	RealCreator   string
	NoteTitle     string
	NoteType      string
	NoteStatus    string
	ReviewAdminID uint
	ReviewPass    bool
	ReviewTime    time.Time
	IfNeedCredit  bool
	CostCredit    int
	NewestVersion uint
	NoteBodies    []NoteBody
}

type NoteCategory struct {
	gorm.Model
	CatID   uint
	CatName string
}

type NoteBody struct {
	gorm.Model
	NoteID          uint
	NoteVersion     uint
	Attachment      []Attachment
	NoteDescription string
}

type Attachment struct {
	gorm.Model
	// 从Image, Zip, PDF, Word, SpecialNote中选择
	Type       string
	FileHash   string
	Nickname   string
	Url        string
	UploaderID uint
}

// DBAttachmentCreate 在数据库中创建附件记录
func DBAttachmentCreate(attachment *Attachment) error {
	result := dao.DB.Create(&attachment)
	return result.Error
}

// DBAttachmentGetByHash 根据文件哈希在数据库中查找附件记录
func DBAttachmentGetByHash(fileHash string) (*Attachment, error) {
	var attachment Attachment
	result := dao.DB.Where("file_hash = ?", fileHash).First(&attachment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &attachment, nil
}

func DBCreateNote(noteInfo *NoteInfo) bool {
	// 这里假设已经存在一个全局的数据库连接对象叫 dao
	if err := dao.DB.Create(noteInfo).Error; err != nil {
		fmt.Println("Error creating note in DB:", err)
		return false
	}
	return true
}
