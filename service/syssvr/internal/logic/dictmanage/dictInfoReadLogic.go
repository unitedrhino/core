package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictInfoReadLogic {
	return &DictInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictInfoReadLogic) DictInfoRead(in *sys.DictInfoReadReq) (*sys.DictInfo, error) {
	var (
		po  *relationDB.SysDictInfo
		err error
	)
	po, err = relationDB.NewDictInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DictInfoFilter{
		ID:   in.Id,
		Code: in.Code,
	})
	return utils.Copy[sys.DictInfo](po), err
}
