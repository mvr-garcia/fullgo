package seed

import (
	"log"

	"github.com/mvr-garcia/fullgo/api/models"
	"gorm.io/gorm"
)

var users = []models.User{
	{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	var err error

	err = db.Migrator().DropTable(&models.Post{}, &models.User{})
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Create(&users).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	for i, post := range posts {
		post.AuthorID = users[i].ID
	}
	err = db.Create(&posts).Error
	if err != nil {
		log.Fatalf("cannot seed posts table: %v", err)
	}
}
