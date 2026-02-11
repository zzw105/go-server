package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-server/config"
	"go-server/model"
)

type ClassificationTree struct {
	ID       int64                `json:"id"`
	Name     string               `json:"name"`
	Children []ClassificationTree `json:"children"`
}

func GetClassificationTree(c *gin.Context) {
	var list []model.Classification
	config.DB.Order("sort ASC").Find(&list)

	tree := buildTree(0, list)
	c.JSON(200, tree)
}

func buildTree(parentID int64, list []model.Classification) []ClassificationTree {
	var tree []ClassificationTree

	for _, item := range list {
		if item.ParentID == parentID {
			node := ClassificationTree{
				ID:   item.ID,
				Name: item.Name,
			}
			node.Children = buildTree(item.ID, list)
			tree = append(tree, node)
		}
	}
	return tree
}

// UpdateClassificationTree 更新整个分类树，替换数据库中的所有数据
func UpdateClassificationTree(c *gin.Context) {
	var tree []ClassificationTree
	if err := c.ShouldBindJSON(&tree); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除所有现有的分类数据
	if err := tx.Exec("DELETE FROM classifications").Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "删除现有数据失败"})
		return
	}

	// 递归插入新的树结构
	if err := insertTreeNodes(tx, tree, 0, 1); err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "插入新数据失败: " + err.Error()})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(500, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(200, gin.H{"message": "分类树更新成功"})
}

// insertTreeNodes 递归插入树节点
func insertTreeNodes(tx *gorm.DB, nodes []ClassificationTree, parentID int64, level int) error {
	for index, node := range nodes {
		classification := model.Classification{
			ID:       node.ID,
			Name:     node.Name,
			ParentID: parentID,
			Level:    level,
			Sort:     index, // 使用数组索引作为排序值
		}

		// 如果ID为0，让数据库自动生成ID
		if node.ID == 0 {
			if err := tx.Omit("id").Create(&classification).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Create(&classification).Error; err != nil {
				return err
			}
		}

		// 递归插入子节点
		if len(node.Children) > 0 {
			if err := insertTreeNodes(tx, node.Children, classification.ID, level+1); err != nil {
				return err
			}
		}
	}
	return nil
}
