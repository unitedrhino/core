package stores

import (
	"gitee.com/i-Things/share/errors"
	"gorm.io/gorm"
	"strings"
)

func ErrFmt(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*errors.CodeError); ok {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.NotFind.WithStack()
	}
	if strings.Contains(err.Error(), "Duplicate entry") {
		return errors.Duplicate.AddDetail(err)
	}
	return errors.Database.AddDetail(err)
}
