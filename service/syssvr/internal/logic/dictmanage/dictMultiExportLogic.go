package dictmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictMultiExportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictMultiExportLogic {
	return &DictMultiExportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictMultiExportLogic) DictMultiExport(in *sys.DictMultiExportReq) (*sys.DictMultiExportResp, error) {
	pos, err := relationDB.NewDictInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DictInfoFilter{Codes: in.DictCodes, WithDetails: true}, nil)
	if err != nil {
		return nil, err
	}
	return &sys.DictMultiExportResp{Dicts: utils.MarshalNoErr(pos)}, nil
}
