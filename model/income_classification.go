package model

type IncomeClassification struct {
	ID       int64  `gorm:"primaryKey;type:bigint" json:"id"`
	Name     string `json:"name"`
	ParentID int64  `gorm:"type:bigint" json:"parent_id"`
	Level    int    `json:"level"`
	Sort     int    `gorm:"type:int;default:0" json:"-"` // 排序字段，不返回给前端
}

// TreeNode 接口实现
func (c IncomeClassification) GetID() int64        { return c.ID }
func (c IncomeClassification) GetName() string     { return c.Name }
func (c IncomeClassification) GetParentID() int64  { return c.ParentID }
func (c *IncomeClassification) SetID(id int64)      { c.ID = id }
func (c *IncomeClassification) SetParentID(id int64) { c.ParentID = id }
func (c *IncomeClassification) SetLevel(level int)  { c.Level = level }
func (c *IncomeClassification) SetSort(sort int)    { c.Sort = sort }
