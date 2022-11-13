package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"foreignKey;not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("required title")
	}
	if p.Content == "" {
		return errors.New("required content")
	}
	if p.AuthorID < 1 {
		return errors.New("required author")
	}
	return nil
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	err := db.Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		err := db.Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return p, nil
}

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	posts := []Post{}
	err := db.Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for _, i := range posts {
			err := db.Model(&User{}).Where("id = ?", i.AuthorID).Take(&i.Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
	err := db.Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err := db.Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return p, nil
}

func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {
	db = db.Model(&Post{}).Where("id = ?", p.ID).Take(&p).Updates(
		Post{
			Title:     p.Title,
			Content:   p.Content,
			UpdatedAt: time.Now(),
		},
	)
	if db.Error != nil {
		return &Post{}, db.Error
	}

	if p.ID != 0 {
		err := db.Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return p, nil
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
