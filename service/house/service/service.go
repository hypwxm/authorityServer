package service

import (
	"worldbar/service/house/model/locationEnumsM"
	"worldbar/service/house/model/locationOptionsM"
)

func CreateEnums(entity *locationEnumsM.WbHouseEnums) (string, error) {
	return entity.Insert()
}

func ModifyEnums(updateQuery *locationEnumsM.UpdateByIDQuery) error {
	return new(locationEnumsM.WbHouseEnums).Update(updateQuery)
}

func EnumsList(query *locationEnumsM.Query) ([]*locationEnumsM.WbHouseEnums, int64, error) {
	return new(locationEnumsM.WbHouseEnums).List(query)
}

func DeleteEnums(query *locationEnumsM.DeleteQuery) error {
	return new(locationEnumsM.WbHouseEnums).Delete(query)
}

func UpdateSort(query *locationEnumsM.UpdateSortQuery) error {
	return new(locationEnumsM.WbHouseEnums).UpdateSort(query)
}

func CreateEnumsOption(entity *locationOptionsM.WbHouseOption) (string, error) {
	return entity.Insert()
}

func ModifyEnumsOption(updateQuery *locationOptionsM.UpdateByIDQuery) error {
	return new(locationOptionsM.WbHouseOption).Update(updateQuery)
}

func EnumsOptionsList(query *locationOptionsM.Query) ([]*locationOptionsM.WbHouseOption, int64, error) {
	return new(locationOptionsM.WbHouseOption).List(query)
}

func DeleteEnumsOption(query *locationOptionsM.DeleteQuery) error {
	return new(locationOptionsM.WbHouseOption).Delete(query)
}

func GetEnumsOption(query *locationOptionsM.GetQuery) (*locationOptionsM.WbHouseOption, error) {
	return new(locationOptionsM.WbHouseOption).GetByID(query)
}

func Associate(query *locationOptionsM.AssociateQuery) error {
	return new(locationOptionsM.WbHouseOptionAssociate).Associate(query)
}

func GetAssociates(query *locationOptionsM.AssociateGetQuery) ([]*locationOptionsM.Associate, error) {
	return new(locationOptionsM.WbHouseOptionAssociate).GetAssociate(query)
}

func DeleteAssociates(query *locationOptionsM.AssociateQuery) error {
	return new(locationOptionsM.WbHouseOptionAssociate).DeleteAssociates(query)
}

func GetOptionsByIds(query *locationOptionsM.GetListByIdsQuery) ([]*locationOptionsM.OptionsWithEnums, error) {
	return new(locationOptionsM.WbHouseOption).GetListByIds(query)
}
