package cache

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/domain/slot"
	"github.com/dgraph-io/ristretto"
	"time"
)

type Slot struct {
	cache *ristretto.Cache
}

const (
	expireTime = time.Minute * 10
)

func NewSlot() *Slot {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
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
