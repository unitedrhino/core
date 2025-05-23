package areamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreaInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoDeleteLogic {
	return &AreaInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除区域
func (l *AreaInfoDeleteLogic) AreaInfoDelete(in *sys.AreaWithID) (*sys.Empty, error) {
	if in.AreaID == 0 {
		return nil, errors.Parameter
	}
	area, err := relationDB.NewAreaInfoRepo(l.ctx).FindOne(l.ctx, in.AreaID, nil)
	if err != nil {
		return nil, err
	}
	if area.DeviceCount > 0 {
		return nil, errors.Parameter.AddDetail(in.AreaID).WithMsg("区域下有设备禁止删除")
	}
	var (
		areaPo *relationDB.SysAreaInfo
	)
	ti, err := relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{})
	if err != nil {
		return nil, err
	}
	if ti.DefaultAreaID == in.AreaID || area.IsSysCreated == def.True {
		return nil, errors.Parameter.AddDetail(in.AreaID).WithMsg("默认区域禁止删除")
	}
	var areaIDs []int64
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		areaPo, err = checkArea(l.ctx, tx, in.AreaID)
		if err != nil {
			return errors.Fmt(err).WithMsg("检查区域出错")
		} else if areaPo == nil {
			return errors.Parameter.AddDetail(in.AreaID).WithMsg("检查区域不存在")
		}
		AiDB := relationDB.NewAreaInfoRepo(tx)

		if areaPo.ParentAreaID != def.RootNode { //如果有父节点
			parent, err := AiDB.FindOne(l.ctx, areaPo.ParentAreaID, nil)
			if err != nil {
				return err
			}
			parent.LowerLevelCount--
			if parent.LowerLevelCount == 0 {
				parent.IsLeaf = def.True
			}
			subSubAreaIDs(l.ctx, tx, parent, in.AreaID)
			err = AiDB.Update(l.ctx, parent)
			if err != nil {
				return err
			}
		}
		areas, err := AiDB.FindByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDPath: areaPo.AreaIDPath}, nil)
		if err != nil {
			return errors.Fmt(err).WithMsg("查询区域及子区域出错")
		}
		for _, area := range areas {
			areaIDs = append(areaIDs, int64(area.AreaID))
		}
		err = relationDB.NewAreaInfoRepo(tx).DeleteByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDs: areaIDs})
		if err != nil {
			return errors.Fmt(err).WithMsg("删除区域及子区域出错")
		}
		err = relationDB.NewDataAreaRepo(tx).DeleteByFilter(l.ctx, relationDB.DataAreaFilter{AreaIDs: areaIDs})
		if err != nil {
			return err
		}
		err = relationDB.NewUserAreaApplyRepo(tx).DeleteByFilter(l.ctx, relationDB.UserAreaApplyFilter{AreaIDs: areaIDs})
		if err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		FillProjectAreaCount(l.ctx, l.svcCtx, int64(areaPo.ProjectID))
		err = l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreAreaInfoDelete, def.IDs{IDs: areaIDs})
		if err != nil {
			l.Error(err)
		}

	}
	l.svcCtx.AreaCache.SetData(l.ctx, in.AreaID, nil)
	return &sys.Empty{}, err
}
