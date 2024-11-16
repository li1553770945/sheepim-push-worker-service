package converter

import (
	"github.com/li1553770945/sheepim-online-service/biz/constant"
	"github.com/li1553770945/sheepim-online-service/biz/internal/domain"
	"github.com/li1553770945/sheepim-online-service/kitex_gen/base"
	"github.com/li1553770945/sheepim-online-service/kitex_gen/project"
)

func AssembleSuccessBaseResp() *base.BaseResp {
	return &base.BaseResp{
		Code: constant.Success,
	}
}

func ProjectInfoEntityToDTO(projectEntities *[]domain.ProjectEntity) *project.ProjectsResp {
	var projectDataList []*project.ProjectData

	// 遍历每个 ProjectEntity，将其转换为 ProjectData
	for _, entity := range *projectEntities {
		// 将 StartDate 转换为时间戳
		startDateTimestamp := entity.StartDate.Unix()

		// 将 EndDate 转换为时间戳，如果为空则保持 nil
		var endDateTimestamp *int64
		if entity.EndDate != nil {
			ts := entity.EndDate.Unix()
			endDateTimestamp = &ts
		}

		// 创建 ProjectData 并赋值
		projectData := &project.ProjectData{
			Name:         entity.Name,
			Desc:         entity.Desc,
			Link:         entity.Link,
			Status:       int32(entity.Status),
			VolumeOfWork: int32(entity.VolumeOfWork),
			Difficulty:   int32(entity.Difficulty),
			StartDate:    startDateTimestamp,
			EndDate:      endDateTimestamp,
		}

		// 将转换后的 projectData 添加到列表中
		projectDataList = append(projectDataList, projectData)
	}

	// 构建并返回 ProjectsResp
	return &project.ProjectsResp{
		BaseResp:    AssembleSuccessBaseResp(),
		ProjectData: projectDataList,
	}
}
