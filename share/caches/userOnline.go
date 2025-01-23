package caches

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/caches"
	"github.com/spf13/cast"
)

type UserOnline struct {
}

func (u *UserOnline) genKey(userID int64) string {
	return fmt.Sprintf("user:online:%v", userID)
}

func (u *UserOnline) SetUser(ctx context.Context, nodeID int64, userID int64) error {
	return caches.GetStore().SetexCtx(ctx, u.genKey(userID), cast.ToString(nodeID), 20)
}
func (u *UserOnline) DelUser(ctx context.Context, userID int64) error {
	_, err := caches.GetStore().DelCtx(ctx, u.genKey(userID))
	return err
}
