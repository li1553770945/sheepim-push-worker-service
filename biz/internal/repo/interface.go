package repo

import (
	"github.com/li1553770945/sheepim-online-service/biz/internal/domain"
	"gorm.io/gorm"
)

type IRepository interface {
	SaveProject(project *domain.ProjectEntity) error
	RemoveProject(projectID int32) error
	GetProject(projectID int32) (*domain.ProjectEntity, error)
	GetProjectsNum() (int64, error)
	GetProjects(from int32, end int32, order string, status int32) (*[]domain.ProjectEntity, error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	err := db.AutoMigrate(&domain.ProjectEntity{})
	if err != nil {
		panic("迁移用户模型失败：" + err.Error())
	}
	return &Repository{
		DB: db,
	}
}
