package detail

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/dict"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典详情单个
func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.DictDetailReadReq) (resp *types.DictDetail, err error) {
	ret, err := l.svcCtx.DictM.DictDetailRead(l.ctx,
		utils.Copy[sys.DictDetailReadReq](req))
	return dict.ToDetailTypes(ret), err
}
