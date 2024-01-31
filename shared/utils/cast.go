package utils

import (
	"database/sql"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
	"time"
)

/*
@in src 赋值的数据源
@in dst 赋值对象的结构体
@out dst类型的结构体
*/
func Convert(src any, dst any) any {

	srcType := reflect.TypeOf(src) //获取type
	dstType := reflect.TypeOf(dst)

	srcEl := reflect.ValueOf(src).Elem() //获取value
	dstEl := reflect.ValueOf(dst).Elem()
	//双循环，对相同名字对字段进行赋值
	for i := 0; i < srcType.NumField(); i++ {
		for j := 0; j < dstType.NumField(); j++ {
			if srcType.Field(i).Name == dstType.Field(j).Name {
				dstEl.Field(i).Set(srcEl.Field(j))
			}
		}
	}
	return dst
}

func ToEmptyInt64(val *wrappers.Int64Value) int64 {
	if val == nil {
		return 0
	}
	return val.Value
}
func ToNullInt64(val *wrappers.Int64Value) *int64 {
	if val == nil {
		return nil
	}
	return &val.Value
}
func ToRpcNullInt64(val any) *wrappers.Int64Value {
	if val == nil {
		return nil
	}

	var wrapVal any
	switch val.(type) {
	case string:
		wrapVal = val
	case *string:
		if v := val.(*string); v != nil {
			wrapVal = *v
		}
	case sql.NullString:
		if v := val.(sql.NullString); v.Valid == true {
			wrapVal = v.String
		}
	case int64:
		wrapVal = val
	case *int64:
		if v := val.(*int64); v != nil {
			wrapVal = *v
		}
	case sql.NullInt64:
		if v := val.(sql.NullInt64); v.Valid == true {
			wrapVal = v.Int64
		}
	default:
		return nil
	}

	if wrapVal != nil {
		return &wrappers.Int64Value{Value: cast.ToInt64(wrapVal)}
	} else {
		return nil
	}
}

func SqlToString(val sql.NullString) string {
	if !val.Valid {
		return ""
	}
	return val.String
}

func ToEmptyString(val *wrappers.StringValue) string {
	if val == nil {
		return ""
	}
	return val.Value
}
func ToNullString(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}
func ToRpcNullString(val any) *wrappers.StringValue {
	if val == nil {
		return nil
	}
	switch val.(type) {
	case string:
		v := val.(string)
		if v == "" {
			return nil
		}
		return &wrappers.StringValue{
			Value: v,
		}
	case *string:
		v := val.(*string)
		if v != nil {
			return &wrappers.StringValue{
				Value: *v,
			}
		}
	case sql.NullString:
		v := val.(sql.NullString)
		if v.Valid == true {
			return &wrappers.StringValue{Value: v.String}
		}
	}
	return nil
}

func ToRpcNullDouble(val *float64) *wrappers.DoubleValue {
	if val != nil {
		return &wrappers.DoubleValue{
			Value: *val,
		}
	}
	return nil
}

var empty = time.Time{}

func Int64ToTimex(in int64) *time.Time {
	if in == 0 {
		return nil
	}
	ret := time.Unix(in, 0)
	return &ret
}

func TimeToInt64(t time.Time) int64 {
	if t == empty {
		return 0
	}
	return t.Unix()
}
func Time2ToInt64(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return TimeToInt64(*t)
}
func SetToSlice[t constraints.Ordered](in map[t]struct{}) (ret []t) {
	for k := range in {
		ret = append(ret, k)
	}
	return
}

func AnyToNullString(in any) sql.NullString {
	if in == nil || IsNil(in) {
		return sql.NullString{}
	}
	switch in.(type) {
	case string, []byte:
		return sql.NullString{
			String: cast.ToString(in),
			Valid:  true,
		}
	case *wrapperspb.StringValue:
		v := in.(*wrapperspb.StringValue)
		if v == nil {
			return sql.NullString{}
		}
		return sql.NullString{
			String: v.Value,
			Valid:  true,
		}
	}
	str, err := json.Marshal(in)
	if err != nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(str), Valid: true}
}
func SqlNullStringToAny(in sql.NullString, ret any) error {
	if in.Valid == false {
		return nil
	}
	err := json.Unmarshal([]byte(in.String), ret)
	return err
}

func SliceTo[retT any](values []string, cov func(any) retT) []retT {
	var ret []retT
	for _, v := range values {
		ret = append(ret, cov(v))
	}
	return ret
}

func TrimNil[a any](in *a) a {
	if in != nil {
		return *in
	}
	var ret a
	return ret
}

// TimeTo24Sec 转换成 24小时的秒单位表示
func TimeTo24Sec(t time.Time) int64 {
	ret := t.Hour() * 60 * 60
	ret += t.Minute() * 60
	ret += t.Second()
	return int64(ret)
}

func ToTimeX(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
func TimeXToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
