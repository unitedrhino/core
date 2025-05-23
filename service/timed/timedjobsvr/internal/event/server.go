package event

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type Server struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewEventServer(ctx context.Context, svcCtx *svc.ServiceContext) *Server {
	return &Server{svcCtx: svcCtx, ctx: ctx, Logger: logx.WithContext(ctx)}
}

func (s *Server) DataClean() error {
	s.Info("start data clean")
	ctxs.GoNewCtx(s.ctx, func(ctx context.Context) {
		err := relationDB.NewJobLogRepo(ctx).DeleteByFilter(ctx, relationDB.TaskLogFilter{CreatedTime: stores.CmpLt(time.Now().Add(-time.Hour * 24 * 3))}) //只保留三天的日志
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	})
	keys, err := s.svcCtx.Redis.KeysCtx(s.ctx, "timed:sql:*:hash:*")
	if err != nil {
		return err
	}
	days := map[string]struct{}{
		time.Now().Add(time.Hour * 24 * time.Duration(-1)).Format("2006-01-02"): {},
		time.Now().Format("2006-01-02"):                                         {},
	}
	for _, key := range keys {
		fields, err := s.svcCtx.Redis.Hkeys(key)
		if err != nil {
			return err
		}
		if len(fields) == 0 { //如果没有使用了
			_, err := s.svcCtx.Redis.Del(key)
			return err
		}
		for _, field := range fields {
			date, _, find := strings.Cut(field, ":")
			if !find { //如果没有找到
				s.svcCtx.Redis.Hdel(key, field)
			}
			if _, ok := days[date]; ok { //在有效期内
				continue
			}
			s.svcCtx.Redis.Hdel(key, field)
		}
	}

	return nil
}
