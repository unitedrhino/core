package oss

import (
	"gitee.com/i-Things/share/conf"
)

func newOssManager(setting conf.OssConf) (sm Handle, err error) {
	OssType := setting.OssType
	switch OssType {
	case "aliyun":
		sm, err = newAliYunOss(conf.AliYunConf{OssConf: setting})
	case "minio":
		sm, err = newMinio(conf.MinioConf{OssConf: setting})
	}
	return sm, err
}
