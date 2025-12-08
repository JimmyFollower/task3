package blog

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// -------------------------- 模型定义（修复后） --------------------------
type User struct {
	ID        int    `gorm:"primaryKey"` // gorm v2推荐primaryKey而非primary_key
	Name      string `gorm:"size:50;not null;unique"`
	PostCount int    `gorm:"default:0"`
	Posts     []Post `gorm:"foreignKey:UserID"` // gorm v2推荐foreignKey
}

type Post struct {
	ID            int       `gorm:"primaryKey"`
	Title         string    `gorm:"size:100;not null"`
	Content       string    `gorm:"type:text"`
	CommentCount  int       `gorm:"default:0"`
	CommentStatus string    `gorm:"size:20;default:有评论"` // 新增评论状态
	UserID        int       `gorm:"not null"`
	User          User      `gorm:"foreignKey:UserID"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	ID      int    `gorm:"primaryKey"`
	Content string `gorm:"type:text"`
	PostID  int    `gorm:"not null"`
	Post    Post   `gorm:"foreignKey:PostID"`
	UserID  int    `gorm:"not null"` // 这里是评论者ID，
}

// 钩子函数
// Post创建前更新用户文章数
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	var user User
	if err := tx.First(&user, p.UserID).Error; err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}

	if err := tx.Model(&User{}).Where("id=?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1")).Error; err != nil {
		log.Printf("更新用户文章数失败: %v", err)
		return fmt.Errorf("更新用户文章数失败: %w", err)
	}

	log.Printf("用户%d的文章数已更新为%d", p.UserID, user.PostCount+1)
	return nil
}

// Comment删除后更新文章评论数和状态
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var post Post
	if err := tx.First(&post, c.PostID).Error; err != nil {
		return fmt.Errorf("获取文章失败: %w", err)
	}

	// 正确更新评论数（字段名小写）
	if err := tx.Model(&Post{}).Where("id=?", c.PostID).
		Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return fmt.Errorf("更新评论数失败: %w", err)
	}

	// 重新查询最新的文章数据
	var latestPost Post
	if err := tx.First(&latestPost, c.PostID).Error; err != nil {
		return fmt.Errorf("重新查询文章失败: %w", err)
	}

	// 评论数为0时更新状态
	if latestPost.CommentCount == 0 {
		if err := tx.Model(&Post{}).Where("id=?", c.PostID).
			Update("comment_status", "无评论").Error; err != nil {
			return fmt.Errorf("更新评论状态失败: %w", err)
		}
		log.Printf("文章%d评论数为0，状态更新为「无评论」", c.PostID)
	}

	return nil
}

// -------------------------- 关联查询-------------------------
// 查询用户所有文章及评论
func getUserPostsAndComments(db *gorm.DB, userID int) ([]Post, error) {
	var posts []Post
	if err := db.Preload("Comments").Where("user_id=?", userID).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("查询用户文章失败: %w", err)
	}
	return posts, nil
}

// 查询评论最多的文章
func getPostWithMostComments(db *gorm.DB) (Post, error) {
	var post Post
	if err := db.Preload("User").Preload("Comments").
		Order("comment_count desc").Limit(1).First(&post).Error; err != nil {
		return post, fmt.Errorf("查询评论最多的文章失败: %w", err)
	}
	return post, nil
}

// -------------------------- 主函数 --------------------------
func Run(db *gorm.DB) {
	// 自动迁移表（同步字段变更）
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		log.Fatal("表迁移失败:", err)
	}

	// 创建用户
	user := User{Name: "meta5"}
	if err := db.Create(&user).Error; err != nil {
		log.Fatal("创建用户失败:", err)
	}
	fmt.Printf("创建用户: %+v\n", user)

	// 创建第一篇文章（触发BeforeCreate钩子）
	post1 := Post{Title: "grom11", Content: "这是一篇关于grom用法的文章1", UserID: user.ID}
	if err := db.Create(&post1).Error; err != nil {
		log.Fatal("创建文章失败:", err)
	}

	// 创建3条评论
	comment1 := Comment{Content: "meta1很棒的文章", PostID: post1.ID, UserID: user.ID}
	comment2 := Comment{Content: "1meta_Sqlx很棒的文章", PostID: post1.ID, UserID: user.ID}
	comment3 := Comment{Content: "3meta_Sqlx学习了1", PostID: post1.ID, UserID: user.ID}
	if err := db.Create(&comment1).Error; err != nil {
		log.Fatal("创建评论失败:", err)
	}
	if err := db.Create(&comment2).Error; err != nil {
		log.Fatal("创建评论失败:", err)
	}
	if err := db.Create(&comment3).Error; err != nil {
		log.Fatal("创建评论失败:", err)
	}
	// 初始化评论数
	if err := db.Model(&Post{}).Where("id=?", post1.ID).Update("comment_count", 3).Error; err != nil {
		log.Fatal("初始化评论数失败:", err)
	}

	// 创建第二篇文章
	post2 := Post{Title: "Gorm进阶教程", Content: "这是一篇关于Gorm进阶用法的文章", UserID: user.ID}
	if err := db.Create(&post2).Error; err != nil {
		log.Fatal("创建文章失败:", err)
	}
	// 创建2条评论
	comment4 := Comment{Content: "很棒！", PostID: post2.ID, UserID: user.ID}
	comment5 := Comment{Content: "学习", PostID: post2.ID, UserID: user.ID}
	db.Create(&comment4)
	db.Create(&comment5)
	db.Model(&Post{}).Where("id=?", post2.ID).Update("comment_count", 2)

	// 查询用户所有文章及评论
	posts, err := getUserPostsAndComments(db, user.ID)
	if err != nil {
		fmt.Printf("查询用户文章失败: %v\n", err)
	} else {
		fmt.Printf("\n用户(%d)的文章及评论: \n", user.ID)
		for _, p := range posts {
			fmt.Printf("UID:%d PID:%d 文章: %s, 评论数: %d\n", p.UserID, p.ID, p.Title, len(p.Comments))
			for _, c := range p.Comments {
				fmt.Printf("UID:%d PID:%d CID:%d  评论: %s\n", c.UserID, c.PostID, c.ID, c.Content)
			}
		}
	}

	// 查询评论最多的文章
	topPost, err := getPostWithMostComments(db)
	if err != nil {
		fmt.Printf("查询评论最多的文章失败: %v\n", err)
	} else {
		fmt.Printf("\n评论最多的文章: %s, 评论数: %d\n", topPost.Title, topPost.CommentCount)
	}

	// 测试删除评论
	var postBeforeDelete Post
	if err := db.First(&postBeforeDelete, post2.ID).Error; err != nil {
		log.Fatal("查询文章失败:", err)
	}
	fmt.Printf("\n删除前 文章 %s 的 评论数: %d, 状态: %s\n",
		postBeforeDelete.Title, postBeforeDelete.CommentCount, postBeforeDelete.CommentStatus)

	// 删除评论（触发AfterDelete钩子）
	db.Delete(&comment4)
	db.Delete(&comment5)

	// 重新查询删除后的文章
	var postAfterDelete Post
	if err := db.First(&postAfterDelete, post2.ID).Error; err != nil {
		log.Fatal("查询删除后的文章失败:", err)
	}
	fmt.Printf("删除后 文章 %s 的 评论数: %d, 状态: %s\n",
		postAfterDelete.Title, postAfterDelete.CommentCount, postAfterDelete.CommentStatus)
}
