package dictmanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictMultiImportLogic {
	return &DictMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictMultiImportLogic) DictMultiImport(in *sys.DictMultiImportReq) (*sys.DictMultiImportResp, error) {
	var pos []*relationDB.SysDictInfo
	err := json.Unmarshal([]byte(in.Dicts), &pos)
	if err != nil {
		return nil, err
	}
	var resp = sys.DictMultiImportResp{Total: int64(len(pos))}
	for _, v := range pos {
		_, err := NewDictInfoCreateLogic(l.ctx, l.svcCtx).DictInfoCreate(utils.Copy[sys.DictInfo](v))
		if err != nil && !errors.Cmp(err, errors.Duplicate) {
			resp.ErrCount++
			l.Error(v, err)
			continue
		}
		resp.SuccCount++
		if len(v.Details) > 0 {
			_, err := NewDictDetailMultiCreateLogic(l.ctx, l.svcCtx).DictDetailMultiCreate(&sys.DictDetailMultiCreateReq{
				DictCode: v.Code,
				List:     ToDictDetailsPb(v.Details),
			})
			if err != nil && !errors.Cmp(err, errors.Duplicate) {
				l.Error(v, err)
			}
		}
	}
	return &resp, nil
}
