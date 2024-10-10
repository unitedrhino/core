package middleware

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/config"
	role "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	user "gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"net/http"
)

type CheckTokenWareMiddleware struct {
	cfg     config.Config
	UserRpc user.UserManage
	AuthRpc role.RoleManage
}

func NewCheckTokenWareMiddleware(cfg config.Config, UserRpc user.UserManage, AuthRpc role.RoleManage) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{cfg: cfg, UserRpc: UserRpc, AuthRpc: AuthRpc}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
