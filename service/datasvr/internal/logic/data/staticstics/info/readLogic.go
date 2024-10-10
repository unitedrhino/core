package info

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/datasvr/internal/domain"
	"gitee.com/unitedrhino/core/service/datasvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/domain/slot"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
	"time"

	"gitee.com/unitedrhino/core/service/datasvr/internal/svc"
	"gitee.com/unitedrhino/core/service/datasvr/internal/types"

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

var (
	ColFmtMap = map[string]string{
		"dayFmt":   "date(%v)",
		"hourFmt":  "DATE_FORMAT(%v, '%%Y-%%m-%%d %%H:00:00')",
		"mouthFmt": "DATE_FORMAT(%v, '%%Y-%%m')",
		"yearFmt":  "Year(%v)",
	}
)

func ColFmt(old string, isToHump bool, as bool) string {
	f, col, found := strings.Cut(old, ":")
	newCol := col
	if col == "" {
		newCol = f
	}
	if isToHump {
		newCol = utils.CamelCaseToUdnderscore(newCol)
	}
	newCol = stores.Col(newCol)
	if !found {
		return newCol
	}
	fu := ColFmtMap[f]
	if fu != "" {
		if as {
			return fmt.Sprintf(fu+" as %v", newCol, col)
		}
		return fmt.Sprintf(fu, newCol)
	}
	return old
}

func FilterFmt(conn *gorm.DB, si *relationDB.DataStatisticsInfo, k string, v any) *gorm.DB {
	f, ok := si.Filter[k]
	if !ok {
		cols := strings.Split(k, ":")
		var fu string
		newCol := cols[0]
		if len(cols) == 3 {
			newCol = ColFmt(strings.Join(cols[:2], ":"), si.IsToHump == def.True, false)
			fu = cols[2]
		} else {
			if len(cols) == 2 {
				fu = cols[1]
			}
			if si.IsToHump == def.True {
				newCol = utils.CamelCaseToUdnderscore(newCol)
			}
			newCol = stores.Col(newCol)
		}

		if len(cols) > 1 {
			switch fu {
			case "jsonEq":
				col1, obj, ok := strings.Cut(cols[0], ".")
				if !ok {
					conn.AddError(errors.Parameter.AddMsg("jsonEq的格式为xxx.xxx:jsonEq"))
					return conn
				}
				if si.IsToHump == def.True {
					obj = utils.CamelCaseToUdnderscore(obj)
					col1 = utils.CamelCaseToUdnderscore(col1)
				}
				col1 = stores.Col(col1)
				conn = conn.Where(fmt.Sprintf("json_extract(%s,'$.%s')=?", col1, obj), v)
				return conn
			case "in":
				v = strings.Split(cast.ToString(v), ",")
				conn = conn.Where(fmt.Sprintf("%s in ?", newCol), v)
				return conn
			case "eq":
				conn = conn.Where(fmt.Sprintf("%s = ?", newCol), v)
				return conn
			case "subChildren":
				if len(cast.ToString(v)) > 0 {
					conn = conn.Where(fmt.Sprintf("%s like ?", newCol), cast.ToString(v)+"%")
				}
				return conn
			default:
				cmp := stores.GetCmp(fu, v)
				if cmp != nil {
					conn = cmp.Where(conn, newCol)
					return conn
				}
			}
		}
		switch val := v.(type) {
		case string:
			conn = conn.Where(fmt.Sprintf("%s like ?", newCol), "%"+val+"%")
		default:
			conn = conn.Where(fmt.Sprintf("%s = ?", newCol), v)
		}
		return conn
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
	return conn
}

func (l *ReadLogic) Handle(req *types.StaticsticsInfoReadReq) (resp *types.StaticsticsInfoReadResp, err error) {
	db := relationDB.NewStatisticsInfoRepo(l.ctx)
	si, err := db.FindOneByFilter(l.ctx, relationDB.StatisticsInfoFilter{Code: req.Code})
	if err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if si.FilterRoles != "" {
		roles := strings.Split(si.FilterRoles, ",")
		rm := utils.SliceToSet(uc.RoleCodes)
		var find bool
		for _, role := range roles {
			if _, ok := rm[role]; ok {
				find = true
				break
			}
		}
		if !find {
			return nil, errors.Permissions
		}
	}
	var (
		columns []string
		groups  []string
	)
	if req.Columns != "" {
		columns = strings.Split(req.Columns, ",")
	}
	if req.GroupBy != "" {
		groups = strings.Split(req.GroupBy, ",")
	}

	for i, c := range columns {
		columns[i] = ColFmt(c, si.IsToHump == def.True, true)
	}
	for i, c := range groups {
		groups[i] = ColFmt(c, si.IsToHump == def.True, false)
	}
	conn := stores.GetTenantConn(l.ctx)
	if si.FilterSlotCode != "" {
		sl, err := l.svcCtx.Slot.GetData(l.ctx, sysExport.GenSlotCacheKey(slot.CodeDataFilter, si.FilterSlotCode))
		if err != nil {
			return nil, err
		}
		var f domain.DateFilterSlotResp
		err = sl.Request(l.ctx, domain.DateFilterSlotReq{Code: si.Code}, &f)
		if err != nil {
			return nil, err
		}
		if f.Where != "" {
			conn = conn.Where(f.Where)
		}
	}
	if si.IsFilterTenant == def.True && uc.TenantCode != def.TenantCodeDefault {
		conn = conn.Where("tenant_code=?", uc.TenantCode)
	}
	if req.Page != nil {
		conn = utils.Copy[stores.PageInfo](req.Page).ToGorm(conn)
	}

	if (req.Page == nil || len(req.Page.Orders) == 0) && si.OrderBy != "" {
		conn = conn.Order(si.OrderBy)
	}
	if si.IsFilterProject == def.True {
		conn = stores.GenProjectAuthScope(l.ctx, conn)
		if si.IsFilterArea == def.True {
			conn = stores.GenAreaAuthScope(l.ctx, conn)
		}
	}
	if si.IsSoftDelete == def.True {
		conn = conn.Where("deleted_time= 0")
	}
	//var columns = "*"
	switch si.Type {
	case "table":
		conn = conn.Table(si.Table)
	case "sql":
		si.Table = "tb"
		conn = conn.Table(fmt.Sprintf("(%s)as tb", si.Sql))
	}

	//填充过滤条件
	for k, v := range req.Filter {
		conn = FilterFmt(conn, si, k, v)
	}

	if len(req.Aggregations) > 0 {
		for _, agg := range req.Aggregations {
			newCol := agg.Column
			if si.IsToHump == def.True {
				newCol = utils.CamelCaseToUdnderscore(newCol)
			}
			if agg.Column == "total" {
				column := fmt.Sprintf("%s(1) as %s", agg.Func, agg.Column)
				columns = append(columns, column)
			} else {
				column := fmt.Sprintf("%s(%s) as %s", agg.Func, newCol, agg.Column)
				columns = append(columns, column)
			}
		}
	}
	conn = conn.Select(strings.Join(columns, ","))
	if len(groups) != 0 {
		conn = conn.Group(strings.Join(groups, ","))
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
			switch val := v.(type) {
			case time.Time:
				r[k] = utils.ToTimeStr(val)
				//r[utils.UderscoreToLowerCamelCase(k)] = val.Unix()
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
