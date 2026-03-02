package service

import (
	"gorm.io/gorm"
)

// TreeNodeReader 树节点读取接口
type TreeNodeReader interface {
	GetID() int64
	GetName() string
	GetParentID() int64
}

// TreeNode 树节点接口，所有分类模型需要实现此接口
type TreeNode interface {
	TreeNodeReader
	SetID(int64)
	SetName(string)
	SetParentID(int64)
	SetLevel(int)
	SetSort(int)
}

// TreeNodeDTO 通用的树节点 DTO
type TreeNodeDTO struct {
	ID       int64         `json:"id" binding:"required"`
	Name     string        `json:"name" binding:"required"`
	Children []TreeNodeDTO `json:"children"`
}

// UpdateTreeRequest 更新树的请求结构
type UpdateTreeRequest struct {
	Tree []TreeNodeDTO `json:"tree" binding:"required"`
}

// BuildTree 构建树结构（通用方法）
func BuildTree[T TreeNodeReader](parentID int64, list []T) []TreeNodeDTO {
	var tree []TreeNodeDTO

	for _, item := range list {
		if item.GetParentID() == parentID {
			node := TreeNodeDTO{
				ID:   item.GetID(),
				Name: item.GetName(),
			}
			node.Children = BuildTree(item.GetID(), list)
			tree = append(tree, node)
		}
	}
	return tree
}

// UpdateTree 更新整个树（通用方法）
func UpdateTree[T TreeNode](tx *gorm.DB, tableName string, tree []TreeNodeDTO, createModel func() T) error {
	// 删除所有现有数据
	if err := tx.Exec("DELETE FROM " + tableName).Error; err != nil {
		return err
	}

	// 递归插入新的树结构
	return insertTreeNodes(tx, tree, createModel, 0, 1)
}

// insertTreeNodes 递归插入树节点
func insertTreeNodes[T TreeNode](tx *gorm.DB, nodes []TreeNodeDTO, createModel func() T, parentID int64, level int) error {
	for index, node := range nodes {
		model := createModel()
		model.SetID(node.ID)
		model.SetName(node.Name)
		model.SetParentID(parentID)
		model.SetLevel(level)
		model.SetSort(index)

		// 如果ID为0，让数据库自动生成ID
		if node.ID == 0 {
			if err := tx.Omit("id").Create(model).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Create(model).Error; err != nil {
				return err
			}
		}

		// 递归插入子节点
		if len(node.Children) > 0 {
			if err := insertTreeNodes(tx, node.Children, createModel, model.GetID(), level+1); err != nil {
				return err
			}
		}
	}
	return nil
}
