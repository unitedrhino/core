package tenant

// 租户信息表
type Info struct {
	ID               int64  `json:"id"`
	Code             string `json:"code"`               // 租户编码
	Name             string `json:"name"`               // 租户名称
	AdminUserID      int64  `json:"adminUserID,string"` // 超级管理员id
	AdminRoleID      int64  `json:"adminRoleID"`        // 超级管理员id
	BackgroundImg    string `json:"backgroundImg"`      //应用首页
	LogoImg          string `json:"logoImg"`            //应用logo地址
	Desc             string `json:"desc"`               //应用描述
	CreatedTime      int64  `json:"createdTime,string"`
	ProjectLimit     int64  `json:"projectLimit"` //项目数量限制: 1: 单项目 2: 多项目
	DefaultProjectID int64  `json:"defaultProjectID"`
	DefaultAreaID    int64  `json:"defaultAreaID"`
}
