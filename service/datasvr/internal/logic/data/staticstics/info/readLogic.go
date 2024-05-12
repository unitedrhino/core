package info

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/datasvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"strings"
	"time"

	"gitee.com/i-Things/core/service/datasvr/internal/svc"
	"gitee.com/i-Things/core/service/datasvr/internal/types"

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

func (l *ReadLogic) Read(req *types.StaticsticsInfoReadReq) (resp *types.StaticsticsInfoReadResp, err error) {
	return l.Handle(req)
}

func (l *ReadLogic) Handle(req *types.StaticsticsInfoReadReq) (resp *types.StaticsticsInfoReadResp, err error) {
	db := relationDB.NewStatisticsInfoRepo(l.ctx)
	si, err := db.FindOneByFilter(l.ctx, relationDB.StatisticsInfoFilter{Code: req.Code})
	if err != nil {
		return nil, err
	}

	if si.IsToHump == def.True {
		if req.GroupBy != "" {
			req.GroupBy = utils.CamelCaseToUdnderscore(req.GroupBy)
		}
		if req.Columns != "" {
			req.Columns = utils.CamelCaseToUdnderscore(req.Columns)
		}
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	conn := stores.GetTenantConn(l.ctx)
	if si.IsFilterTenant == def.True {
		conn = conn.Where("tenant_code=?", uc.TenantCode)
	}
	if si.IsFilterProject == def.True {
		conn = conn.Where("project_id=?", uc.ProjectID)
	}
	if si.IsSoftDelete == def.True {
		conn = conn.Where("deleted_time= 0")
	}
	if si.IsFilterArea == def.True {
		//todo
	}
	//var columns = "*"
	switch si.Type {
	case "table":
		conn = conn.Table(si.Table)
	case "sql":
		si.Table = "tb"
		conn = conn.Table(fmt.Sprintf("(%s)as tb", si.Sql))
	}
	if si.OrderBy != "" {
		conn = conn.Order(si.OrderBy)
	}
	//填充过滤条件
	for k, v := range req.Filter {
		f, ok := si.Filter[k]
		if !ok {
			switch v.(type) {
			case string:
				conn = conn.Where(fmt.Sprintf("%s like ?", utils.CamelCaseToUdnderscore(k)), "%"+v.(string)+"%")
			default:
				conn = conn.Where(fmt.Sprintf("%s = ?", utils.CamelCaseToUdnderscore(k)), v)
			}
			//return nil, errors.Parameter.WithMsgf("过滤的key未定义:%s", k)
		} else {
			var args []interface{}
			for i := int64(0); i < f.ValNum; i++ {
				switch f.Type {
				case "date":
					v = utils.FmtDateStr(cast.ToString(v))
				case "array": //数组类型
					v = strings.Split(cast.ToString(v), ",")
				}
				args = append(args, v)
			}
			conn = conn.Where(f.Sql, args...)
		}

	}
	var (
		columns []string
	)
	if req.Columns != "" {
		columns = strings.Split(req.Columns, ",")
	}
	if len(req.Aggregations) > 0 {

		//column := req.Columns
		//if column == "" {
		//	column = fmt.Sprintf("%s.*", si.Table)
		//} else {
		//	columns := strings.Split(column, ",")
		//	for i, v := range columns {
		//		columns[i] = fmt.Sprintf("%s.%s", si.Table, v)
		//	}
		//	column = strings.Join(columns, ",")
		//}
		for _, agg := range req.Aggregations {
			as := agg.AsName
			if as == "" {
				as = agg.Column
			}
			if val := si.ArgColumns[agg.Column]; val != "" {
				column := fmt.Sprintf("%s(%s) as %s", agg.Func, val, as)
				columns = append(columns, column)
			} else if agg.Column == "total" {
				column := fmt.Sprintf("%s(1) as %s", agg.Func, as)
				columns = append(columns, column)
			} else {
				column := fmt.Sprintf("%s(%s) as %s", agg.Func, agg.Column, agg.Column)
				columns = append(columns, column)
			}
		}
		//if req.GroupBy != "" {
		//	conn = conn.Select(fmt.Sprintf("%s(%s) as %s,%s",
		//		req.ArgFunc, utils.CamelCaseToUdnderscore(req.ArgColumn), req.ArgFunc, column))
		//} else {
		//	conn = conn.Select(fmt.Sprintf("%s(%s) as %s,%s",
		//		req.ArgFunc, utils.CamelCaseToUdnderscore(req.ArgColumn), req.ArgFunc, column))
		//}
	}
	conn = conn.Select(strings.Join(columns, ","))
	if req.GroupBy != "" {
		conn = conn.Group(utils.CamelCaseToUdnderscore(req.GroupBy))
	}
	var ret = []map[string]any{}
	err = conn.Find(&ret).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	var omitSet = map[string]struct{}{}
	if si.Omits != "" {
		omits := strings.Split(si.Omits, ",")
		for _, v := range omits {
			omitSet[v] = struct{}{}
		}
	}
	for _, r := range ret {
		for k, v := range r {
			delete(r, k)
			if _, ok := omitSet[k]; ok { //忽略的字段
				continue
			}
			switch v.(type) {
			case time.Time:
				r[utils.UderscoreToLowerCamelCase(k)] = v.(time.Time).Unix()
			default:
				k2 := utils.UderscoreToLowerCamelCase(k)
				if utils.SliceIn(k2, "areaID", "projectID") {
					r[k2] = cast.ToString(v)
				} else {
					r[k2] = v
				}
			}
		}
	}
	return &types.StaticsticsInfoReadResp{List: ret}, nil
}
