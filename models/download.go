package models

// 用于记录下载相关信息
type DownloadRecord struct {
	Email string `gorm:"type:varchar(100) not null"`
	ID string `gorm:"type:varchar(32) not null;unique_index"`
	// false 未开始 true 已开始
	Start bool `gorm:"type:tinyint(1) not null;default:0"`
}

// 新建下载记录
func NewDownloadRecord() *DownloadRecord {
	return &DownloadRecord{}
}

func GetRecordByID(ID string) (*DownloadRecord, error) {
	var record DownloadRecord
	result := DB.Where("id = ?", ID).Find(&record)
	return &record, result.Error
}

func GetRecordByEmail(email string) (*DownloadRecord, error) {
	var record DownloadRecord
	result := DB.Where("email = ?", email).Find(&record)
	return &record, result.Error
}

func (record *DownloadRecord) SetStart(start bool) error {
	return DB.Model(&record).Update("start", start).Error
}

// 通过ID删除下载记录
func DeleteRecordByID(ID string) error {
	result := DB.Unscoped().Where("id = ?", ID).Delete(&DownloadRecord{})
	return result.Error
}

