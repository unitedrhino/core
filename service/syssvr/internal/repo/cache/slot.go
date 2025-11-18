package cache

import (
	"context"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/share/domain/slot"
	"github.com/dgraph-io/ristretto"
)

type Slot struct {
	cache *ristretto.Cache
}

const (
	expireTime = time.Minute * 10
)

func NewSlot() *Slot {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,              // number of keys to track frequency of (10M).
		MaxCost:     128 * 1024 * 1024, // 128MB
		BufferItems: 64,                // number of keys per Get buffer.
	})

	return &Slot{
		cache: cache,
	}
}

func (c *Slot) Get(ctx context.Context, code string, subCode string) slot.Infos {
	key := code + ":" + subCode
	v, ok := c.cache.Get(key)
	if ok {
		return v.(slot.Infos)
	}
	list, err := relationDB.NewSlotInfoRepo(ctx).FindByFilter(ctx, relationDB.SlotInfoFilter{Code: code, SubCode: subCode}, nil)
	if err != nil {
		return nil
	}
	slots := relationDB.ToSlotsDo(list)
	c.cache.SetWithTTL(key, slots, 1, expireTime)
	return slots
}
