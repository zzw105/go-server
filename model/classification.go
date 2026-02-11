package model

type Classification struct {
	ID       int64  `gorm:"primaryKey;type:bigint" json:"id"`
	Name     string `json:"name"`
	ParentID int64  `gorm:"type:bigint" json:"parent_id"`
	Level    int    `json:"level"`
	Sort     int    `gorm:"type:int;default:0" json:"-"` // 排序字段，不返回给前端
}
