// 用户管理数据组装逻辑。
package usermanagelogic

import (
	"context"
	"net/url"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

// UserInfoToPb 将用户数据库模型转换为 RPC 用户信息，并为本地私有头像生成临时访问地址。
func UserInfoToPb(ctx context.Context, ui *relationDB.SysUserInfo, svcCtx *svc.ServiceContext) *sys.UserInfo {
	if ui.HeadImg != "" && !isExternalHeadImgURL(ui.HeadImg) {
		var err error
		ui.HeadImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, ui.HeadImg, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return utils.Copy[sys.UserInfo](ui)
}

// isExternalHeadImgURL 判断头像是否已经是第三方可访问地址，避免再次包装成本地 OSS 下载地址。
func isExternalHeadImgURL(headImg string) bool {
	u, err := url.Parse(headImg)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// ToUserAreaApplyInfos 将用户区域申请模型列表转换为 RPC 列表。
func ToUserAreaApplyInfos(in []*relationDB.SysUserAreaApply) (ret []*sys.UserAreaApplyInfo) {
	for _, v := range in {
		ret = append(ret, &sys.UserAreaApplyInfo{
			Id:          v.ID,
			UserID:      v.UserID,
			AreaID:      int64(v.AreaID),
			AuthType:    v.AuthType,
			CreatedTime: v.CreatedTime.Unix(),
		})
	}
	return
}
