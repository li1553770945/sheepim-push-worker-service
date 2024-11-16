package project

import (
	"context"
	"github.com/li1553770945/sheepim-online-service/biz/constant"
	"github.com/li1553770945/sheepim-online-service/biz/internal/converter"
	"github.com/li1553770945/sheepim-online-service/kitex_gen/base"
	"github.com/li1553770945/sheepim-online-service/kitex_gen/project"
)

func (s *ProjectService) GetProjectNum(ctx context.Context) (*project.ProjectNumResp, error) {
	count, err := s.Repo.GetProjectsNum()
	if err != nil {
		return nil, err
	}
	return &project.ProjectNumResp{
		BaseResp: &base.BaseResp{
			Code: constant.Success,
		},
		Num: &count,
	}, nil
}

func (s *ProjectService) GetProjects(ctx context.Context, req *project.ProjectsReq) (*project.ProjectsResp, error) {
	defaultSort := "default"
	var defaultStatus int32 = 0
	var defaultStart int32 = 0
	sort := req.Order
	if sort == nil {
		sort = &defaultSort
	}
	sortFields := map[string]string{
		"default":             "start_date desc",
		"volume_of_work_asc":  "volume_of_work asc",
		"volume_of_work_desc": "volume_of_work desc",
		"start_date_asc":      "start_date asc",
		"start_date_desc":     "start_date desc",
		"difficulty_asc":      "difficulty asc",
		"difficulty_desc":     "difficulty desc",
	}

	order, ok := sortFields[*sort]
	if !ok {
		return &project.ProjectsResp{
			BaseResp: &base.BaseResp{
				Code:    constant.InvalidInput,
				Message: "排序参数错误",
			},
		}, nil

	}

	status := req.Status
	if status == nil {
		status = &defaultStatus
	}

	start := req.Start
	if start == nil {
		start = &defaultStart
	}

	var defaultEnd = *start + 10
	end := req.End
	if end == nil {
		end = &defaultEnd
	}
	if *end-*start > constant.MaxProjectsNum {
		*end = *start + constant.MaxProjectsNum + 1
	}

	projectsResult, err := s.Repo.GetProjects(*start, *end, order, *status)

	if err != nil {
		return &project.ProjectsResp{BaseResp: &base.BaseResp{
			Code:    constant.SystemError,
			Message: "数据库查询失败:" + err.Error(),
		}}, nil
	}

	return converter.ProjectInfoEntityToDTO(projectsResult), nil

}
