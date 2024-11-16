package repo

import "github.com/li1553770945/sheepim-online-service/biz/internal/domain"

func (Repo *Repository) RemoveProject(projectID int32) error {
	err := Repo.DB.Delete(&domain.ProjectEntity{}, projectID).Error
	if err != nil {
		return err
	}
	return nil
}

func (Repo *Repository) SaveProject(project *domain.ProjectEntity) error {
	if project.ID == 0 {
		err := Repo.DB.Create(&project).Error
		return err
	} else {
		err := Repo.DB.Save(&project).Error
		return err
	}
}

func (Repo *Repository) GetProject(projectID int32) (*domain.ProjectEntity, error) {
	var project domain.ProjectEntity
	err := Repo.DB.Where("id = ?", projectID).Limit(1).Find(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}
func (Repo *Repository) GetProjectsNum() (int64, error) {
	var num int64
	err := Repo.DB.Model(&domain.ProjectEntity{}).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}
func (Repo *Repository) GetProjects(start int32, end int32, order string, status int32) (*[]domain.ProjectEntity, error) {
	var projects []domain.ProjectEntity
	db := Repo.DB
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	err := db.Order(order).Order("created_at desc").Offset(int(start)).Limit(int(end - start)).Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return &projects, nil
}
