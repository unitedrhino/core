package common

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadFileLogic) UploadFile() (resp *types.UploadFileResp, err error) {
	file, fh, err := l.r.FormFile("file")
	if err != nil {
		return resp, err
	}
	defer file.Close()
	newFilePath, err := oss.GetFilePath2(l.ctx, fh)
	if err != nil {
		l.Errorf("%s.GetFilePath err:%v", utils.FuncName(), err)
		return nil, err
	}
	if isForbiddenExtension(newFilePath) {
		return nil, errors.Parameter.AddDetail(fmt.Sprintf("有客户上传危险文件 文件名:%s uc:%#v", newFilePath, ctxs.GetUserCtx(l.ctx)))
	}
	fileUri, err := l.svcCtx.OssClient.TemporaryBucket().Upload(l.ctx, newFilePath, file, common.OptionKv{})
	if err != nil {
		return resp, err
	}
	return &types.UploadFileResp{
		FileUri:  fileUri,
		FilePath: newFilePath,
	}, err
}

// 定义一个禁止上传的文件后缀集合
var forbiddenExtensions = map[string]struct{}{
	"html": {}, "htm": {},
	"php": {}, "php5": {}, "php4": {}, "php3": {}, "php2": {}, "phtml": {}, "pht": {},
	"asp": {}, "aspx": {}, "asa": {}, "asax": {}, "ascx": {}, "ashx": {}, "asmx": {}, "cer": {},
	"jsp": {}, "jspa": {}, "jspx": {}, "jsw": {}, "jsv": {}, "jspf": {}, "jhtml": {},
	"htaccess": {}, "swf": {},
}

// 检查文件后缀是否被禁止
func isForbiddenExtension(filename string) bool {
	// 获取文件后缀
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	// 检查是否在禁止集合中
	_, exists := forbiddenExtensions[ext]
	return exists
}
