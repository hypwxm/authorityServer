package authorityServer // import "github.com/hypwxm/authorityServer"

import (
	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/config"
	"github.com/hypwxm/authorityServer/router"
	"github.com/hypwxm/authorityServer/service/admin/menu/dbModel"
	dbModel1 "github.com/hypwxm/authorityServer/service/admin/org/dbModel"
	dbModel2 "github.com/hypwxm/authorityServer/service/admin/role/dbModel"
	dbModel3 "github.com/hypwxm/authorityServer/service/admin/rolePermission/menu/dbModel"
	dbModel4 "github.com/hypwxm/authorityServer/service/admin/user/dbModel"

	"github.com/hypwxm/rider"
)

func Init(app *rider.Rider, cfg config.ConfigDefine) *rider.Router {
	config.InitConfig(cfg)
	// app.Kid("/", router.Router())
	return router.Router()
}

func InitTables() {
	db := appGorm.Open()
	db.AutoMigrate(&dbModel.GMenu{})
	db.AutoMigrate(&dbModel1.GOrg{})
	db.AutoMigrate(&dbModel2.GRole{})
	db.AutoMigrate(&dbModel3.GRoleMenu{})
	db.AutoMigrate(&dbModel4.GAdminUser{})
	db.AutoMigrate(&dbModel4.GUserRole{})
}
