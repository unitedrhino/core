package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoIndexLogic {
	return &ModuleInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoIndexLogic) ModuleInfoIndex(in *sys.ModuleInfoIndexReq) (*sys.ModuleInfoIndexResp, error) {
	f := relationDB.ModuleInfoFilter{Codes: in.Codes, Code: in.Code, Name: in.Name}
	if in.AppCode != "" {
		am, err := relationDB.NewAppModuleRepo(l.ctx).FindByFilter(l.ctx, relationDB.AppModuleFilter{AppCodes: []string{in.AppCode}}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range am {
			f.Codes = append(f.Codes, v.ModuleCode)
		}
	}
	ret, err := relationDB.NewModuleInfoRepo(l.ctx).FindByFilter(l.ctx,
		f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewModuleInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &sys.ModuleInfoIndexResp{List: logic.ToModuleInfosPb(ret), Total: total}, nil
}
