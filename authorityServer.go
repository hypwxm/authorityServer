package authorityServer // import "github.com/hypwxm/authorityServer"

import (
	"log"

	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/config"
	"github.com/hypwxm/authorityServer/router"
	"github.com/hypwxm/authorityServer/service/admin/menu/dbModel"
	dbModel1 "github.com/hypwxm/authorityServer/service/admin/org/dbModel"
	"github.com/hypwxm/authorityServer/service/admin/org/service"
	userService "github.com/hypwxm/authorityServer/service/admin/user/service"

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
	db.Migrator().DropTable(&dbModel.GMenu{})
	db.AutoMigrate(&dbModel.GMenu{})

	db.Migrator().DropTable(&dbModel1.GOrg{})
	db.AutoMigrate(&dbModel1.GOrg{})

	db.Migrator().DropTable(&dbModel2.GRole{})
	db.AutoMigrate(&dbModel2.GRole{})

	db.Migrator().DropTable(&dbModel3.GRoleMenu{})
	db.AutoMigrate(&dbModel3.GRoleMenu{})

	db.Migrator().DropTable(&dbModel4.GAdminUser{})
	db.AutoMigrate(&dbModel4.GAdminUser{})

	db.Migrator().DropTable(&dbModel4.GUserRole{})
	db.AutoMigrate(&dbModel4.GUserRole{})
}

func InitData() {
	org := &dbModel1.GOrg{
		Name: "顶级组织",
	}
	_, err := service.Create(&service.CreateModel{GOrg: *org})
	if err != nil {
		log.Fatal(err)
	}

	err = userService.InitAdmin()
	if err != nil {
		log.Fatal(err)
	}
}
