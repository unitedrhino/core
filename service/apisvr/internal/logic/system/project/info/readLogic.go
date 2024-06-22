package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ProjectWithID) (resp *types.ProjectInfo, err error) {
	dmResp, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: req.ProjectID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ProjectManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	user, err := l.svcCtx.ProjectRpc.ProjectInfoRead(ctxs.WithRoot(l.ctx), &sys.ProjectInfoReadReq{
		ProjectID: dmResp.AdminProjectID,
	})

	return system.ProjectInfoToApi(dmResp, user), nil
}
