package model

type Classification struct {
	ID       int64  `gorm:"primaryKey;type:bigint" json:"id"`
	Name     string `json:"name"`
	ParentID int64  `gorm:"type:bigint" json:"parent_id"`
	Level    int    `json:"level"`
	Sort     int    `gorm:"type:int;default:0" json:"-"` // 排序字段，不返回给前端
}

// TreeNode 接口实现
func (c Classification) GetID() int64        { return c.ID }
func (c Classification) GetName() string     { return c.Name }
func (c Classification) GetParentID() int64  { return c.ParentID }
func (c *Classification) SetID(id int64)        { c.ID = id }
func (c *Classification) SetName(name string)   { c.Name = name }
func (c *Classification) SetParentID(id int64)  { c.ParentID = id }
func (c *Classification) SetLevel(level int)  { c.Level = level }
func (c *Classification) SetSort(sort int)    { c.Sort = sort }
