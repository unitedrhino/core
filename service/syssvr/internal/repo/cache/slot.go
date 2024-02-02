package cache

import (
	"github.com/zeromicro/go-zero/core/stores/kv"
)

type Slot struct {
	store kv.Store
}

func NewSlot(store kv.Store) *Slot {
	return &Slot{
		store: store,
	}
}

//func (c *Slot) GenKey(code string) string {
//	return "slot:" + code
//}
//func (c *Slot) Update(ctx context.Context, code string) error {
//	key := c.GenKey(code)
//	list, err := relationDB.NewSlotInfoRepo(ctx).FindByFilter(ctx, relationDB.SlotInfoFilter{Code: code}, nil)
//	if err != nil {
//		return err
//	}
//	slots := relationDB.ToSlotsDo(list)
//	val, err := c.store.GetCtx(ctx, key)
//	if err != nil || val == "" {
//		return ""
//	}
//	//如果验证码存在，则删除验证码
//	c.store.DelCtx(ctx, key)
//	body := map[string]string{}
//	json.Unmarshal([]byte(val), &body)
//	if body["code"] == code {
//		if body["account"] == "" {
//			return " "
//		}
//		return body["account"]
//	}
//	return ""
//}
//
//func (c *Slot) Get(ctx context.Context, code string) error {
//	body := map[string]interface{}{
//		"code":    code,
//		"account": account,
//	}
//	bodytStr, _ := json.Marshal(body)
//	return c.store.SetexCtx(ctx, c.GenKey(Type, Use, codeID), string(bodytStr), int(expire))
//}
