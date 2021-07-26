package structer

import "github.com/hypwxm/authorityServer/util/database"

type Media struct {
	database.BaseColumns
	Url      string `db:"url" json:"url"`
	Creator  string `db:"creator" json:"creator"`
	Business string `db:"business" json:"business"`
}

func (self *Media) Init() {
	self.BaseColumns.Init()
}

func InitMedias(list []*Media, businessName string, creator string) []*Media {
	for _, v := range list {
		v.Init()
		v.Business = businessName
		v.Creator = creator
	}
	return list
}
