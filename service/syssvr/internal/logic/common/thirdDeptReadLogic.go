package commonlogic

import (
	"context"
	"gitee.com/unitedrhino/share/clients/dingClient"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zhaoyunxing92/dingtalk/v2/request"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdDeptReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThirdDeptReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdDeptReadLogic {
	return &ThirdDeptReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThirdDeptReadLogic) ThirdDeptRead(in *sys.ThirdDeptInfoReadReq) (*sys.DeptInfo, error) {
	cli, err := dingClient.NewDingTalkClient(&conf.ThirdConf{
		AppID:     in.ThirdConfig.AppID,
		AppKey:    in.ThirdConfig.AppKey,
		AppSecret: in.ThirdConfig.AppSecret,
	})
	if err != nil {
		return nil, err
	}

	req := request.NewDeptDetail(int(in.Id))
	detail, err := cli.GetDeptDetail(req.Build())
	if err != nil {
		return nil, errors.Default.WithMsg(err.Error())
	}
	ret := &sys.DeptInfo{
		Id:       int64(detail.Detail.Id),
		Name:     detail.Detail.Name,
		ParentID: int64(detail.Detail.ParentId),
		Desc:     utils.ToRpcNullString(detail.Detail.Brief),
	}
	if in.WithFather && ret.ParentID > def.RootNode {
		req := request.NewDeptDetail(int(ret.ParentID))
		father, err := cli.GetDeptDetail(req.Build())
		if err != nil {
			return nil, errors.Default.WithMsg(err.Error())
		}
		ret.Parent = &sys.DeptInfo{
			Id:       int64(father.Detail.Id),
			Name:     father.Detail.Name,
			ParentID: int64(father.Detail.ParentId),
			Desc:     utils.ToRpcNullString(father.Detail.Brief),
		}
	}
	if in.WithChildren {
		req := request.NewDeptList()
		req.SetDeptId(int(ret.Id))
		dings, err := cli.GetDeptList(req.Build())
		if err != nil {
			return nil, errors.Default.WithMsg(err.Error())
		}
		var list []*sys.DeptInfo
		for _, v := range dings.List {
			list = append(list, &sys.DeptInfo{
				Id:       int64(v.Id),
				Name:     v.Name,
				ParentID: int64(v.ParentId),
			})
		}
		ret.Children = list
	}
	return ret, nil
}
