package tenant

type RegisterAutoCreateProject struct {
	ID           int64                             `json:"id"`
	ProjectName  string                            `json:"projectName"`
	IsSysCreated int64                             `json:"isSysCreated"` //是否是系统创建的,系统创建的只有管理员可以删除
	Areas        []*RegisterAutoCreateArea         `json:"areas"`
	AreaMap      map[int64]*RegisterAutoCreateArea `json:"-"`
	MaxAreaID    int64                             `json:"-"`
}
type RegisterAutoCreateArea struct {
	ID           int64  `json:"id"`
	AreaName     string `json:"areaName"`
	AreaImg      string `json:"areaImg"`
	IsSysCreated int64  `json:"isSysCreated"` //是否是系统创建的,系统创建的只有管理员可以删除
}

func RegisterAutoCreateProjectToMap(in []*RegisterAutoCreateProject) (map[int64]*RegisterAutoCreateProject, int64) {
	var ret = make(map[int64]*RegisterAutoCreateProject)
	var maxID = int64(0)
	for i, v := range in {
		for _, v2 := range v.Areas {
			if in[i].AreaMap == nil {
				in[i].AreaMap = make(map[int64]*RegisterAutoCreateArea)
			}
			in[i].AreaMap[v2.ID] = v2
			if v.MaxAreaID < v2.ID {
				v.MaxAreaID = v2.ID
			}
		}
		ret[v.ID] = v
		if v.ID > maxID {
			maxID = v.ID
		}
	}
	return ret, maxID
}
