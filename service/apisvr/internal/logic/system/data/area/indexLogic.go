package area

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.DataAreaIndexReq) (resp *types.DataAreaIndexResp, err error) {
	ret, err := l.svcCtx.DataM.DataAreaIndex(l.ctx, utils.Copy[sys.DataAreaIndexReq](req))
	if err != nil {
		l.Errorf("%s.rpc.DataAreaIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	if len(ret.List) == 0 {
		return &types.DataAreaIndexResp{}, nil
	}
	var areaIDs []int64
	for _, v := range ret.List {
		areaIDs = append(areaIDs, v.AreaID)
	}
	areaInfos, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{AreaIDs: areaIDs})
	if err != nil {
		return nil, err
	}
	var areaMap = map[int64]*sys.AreaInfo{}
	for _, v := range areaInfos.List {
		areaMap[v.AreaID] = v
	}
	list := ToDataAreaDetail(l.ctx, l.svcCtx, ret.List, areaMap)
	return &types.DataAreaIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     list,
	}, nil
}
