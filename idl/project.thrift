namespace go project
include "base.thrift"

struct ProjectsReq{
    1: optional i32 start
    2: optional i32 end
    3: optional string order
    4: optional i32 status
}
struct ProjectData{
     1: string name,
     2: string desc,
     3: string link,
     4: i32 status,
     5: i32 volume_of_work,
     6: i32 difficulty,
     7: i64 start_date,      // Represented as a timestamp in milliseconds
     8: optional i64 end_date // Represented as a timestamp in milliseconds, optional
   }

struct ProjectsResp {
    1: required base.BaseResp baseResp
    2: optional list<ProjectData> projectData
}

struct ProjectNumResp{
    1: required base.BaseResp baseResp
    2: optional i64 num
}
service ProjectService {
    ProjectsResp GetProjects(ProjectsReq req)
    ProjectNumResp GetProjectNum()
}
