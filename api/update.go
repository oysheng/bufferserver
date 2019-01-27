package api

import (
	"github.com/bytom/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/database/orm"
)

func (s *Server) UpdateBase(c *gin.Context, req *common.AssetProgram) error {
	base := &orm.Base{AssetID: req.Asset, ControlProgram: req.Program}
	if err := s.db.Master().Where(base).First(base).Error; err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "db query base")
	} else if err == gorm.ErrRecordNotFound {
		if err := s.db.Master().Save(base).Error; err != nil {
			return errors.Wrap(err, "update base")
		}
	}

	return nil
}
