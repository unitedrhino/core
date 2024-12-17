package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdDeptIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取第三方的部门信息
func NewThirdDeptIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdDeptIndexLogic {
	return &ThirdDeptIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdDeptIndexLogic) ThirdDeptIndex(req *types.ThirdDeptInfoIndexReq) (resp *types.DeptInfoIndexResp, err error) {
	ret, err := l.svcCtx.Common.ThirdDeptIndex(l.ctx, utils.Copy[sys.ThirdDeptInfoIndexReq](req))

	return utils.Copy[types.DeptInfoIndexResp](ret), err
}
