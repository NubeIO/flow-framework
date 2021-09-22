package database

import (
	"errors"

	"github.com/NubeDev/flow-framework/floweng/core"
	"github.com/NubeDev/flow-framework/model"
	"gorm.io/gorm"
)

func (d *GormDatabase) CreateModel(model interface{}) error {
	query := d.DB.Create(model)
	return query.Error
}

func (d *GormDatabase) GetModel(id int, model *interface{}) error {
	query := d.DB.Where("id = ? ", id).First(&model)
	return query.Error
}

func (d *GormDatabase) GetModelList(modelList interface{}) error {
	query := d.DB.Find(modelList)
	return query.Error
}

func (d *GormDatabase) DeleteModel(id int, model interface{}) (bool, error) {
	query := d.DB.Where("id = ? ", id).Delete(model)
	return query.RowsAffected > 0, query.Error
}

func (d *GormDatabase) UpdateModel(id int, model interface{}) error {
	query := d.DB.Model(model).Where("id = ?", id).Updates(model)
	return query.Error
}

func (d *GormDatabase) DropBlockList(model interface{}) (bool, error) {
	query := d.DB.Where("1 = 1").Delete(model)
	return query.RowsAffected > 0, query.Error
}

func (d *GormDatabase) UpdateBlockName(blockId int, name string) error {
	query := d.DB.Model(&model.Block{}).Where("id = ?", blockId).Update("label", name)
	return query.Error
}

func (d *GormDatabase) UpdateBlockPosition(blockId int, x float64, y float64) error {
	query := d.DB.Model(&model.Block{}).Where("id = ?", blockId).Updates(model.Block{Position: model.Position{X: x, Y: y}})
	return query.Error
}

func (d *GormDatabase) UpdateBlockStaticInput(blockId int, routeIndex int, value *core.InputValue, valueType core.JSONType) error {
	var route *model.BlockStaticRoute
	query := d.DB.Model(&model.BlockStaticRoute{}).Where("block_id = ?", blockId).Where("route_index = ?", routeIndex).First(&route)
	if query.Error == gorm.ErrRecordNotFound {
		route = &model.BlockStaticRoute{
			BlockID:    blockId,
			RouteIndex: routeIndex,
			Type:       int(valueType),
		}
		d.CreateModel(&route)
		var err error = nil
		switch valueType {
		case core.NUMBER:
			err = d.CreateModel(&model.BlockRouteValueNumber{
				BlockRoute: route.ID,
				Value:      value.Data.(float64),
			})
		case core.STRING:
			err = d.CreateModel(&model.BlockRouteValueString{
				BlockRoute: route.ID,
				Value:      value.Data.(string),
			})
		case core.BOOLEAN:
			err = d.CreateModel(&model.BlockRouteValueBool{
				BlockRoute: route.ID,
				Value:      value.Data.(bool),
			})
		default:
			return errors.New("block route type not supported")
		}
		return err
	} else {
		switch valueType {
		case core.NUMBER:
			query = d.DB.Model(&model.BlockRouteValueNumber{}).
				Where("block_route = ?", route.ID).
				Update("value", value.Data)
		case core.STRING:
			query = d.DB.Model(&model.BlockRouteValueString{}).
				Where("block_route = ?", route.ID).
				Update("value", value.Data)
		case core.BOOLEAN:
			query = d.DB.Model(&model.BlockRouteValueBool{}).
				Where("block_route = ?", route.ID).
				Update("value", value.Data)
		default:
			return errors.New("block route type not supported")
		}
	}
	return query.Error
}

func (d *GormDatabase) GetBlockStaticInputs() ([]*model.ProtoBlockStaticRoute, error) {
	var blockRoutes []*model.ProtoBlockStaticRoute

	var blockRoutesNumbers []*model.ProtoBlockRouteNumber
	d.DB.Table("block_route_value_numbers").
		Select("block_static_routes.block_id, block_static_routes.route_index, block_static_routes.type, block_route_value_numbers.value").
		Joins("LEFT JOIN block_static_routes on block_route_value_numbers.block_route = block_static_routes.id").
		Scan(&blockRoutesNumbers)
	for _, route := range blockRoutesNumbers {
		blockRoutes = append(blockRoutes, &model.ProtoBlockStaticRoute{
			BlockStaticRoute: route.BlockStaticRoute,
			Value:            route.Value,
		})
	}

	var blockRoutesStrings []*model.ProtoBlockRouteString
	d.DB.Table("block_route_value_strings").
		Select("block_static_routes.block_id, block_static_routes.route_index, block_static_routes.type, block_route_value_strings.value").
		Joins("LEFT JOIN block_static_routes on block_route_value_strings.block_route = block_static_routes.id").
		Scan(&blockRoutesStrings)
	for _, route := range blockRoutesStrings {
		blockRoutes = append(blockRoutes, &model.ProtoBlockStaticRoute{
			BlockStaticRoute: route.BlockStaticRoute,
			Value:            route.Value,
		})
	}

	var blockRoutesBools []*model.ProtoBlockRouteBool
	d.DB.Table("block_route_value_bools").
		Select("block_static_routes.block_id, block_static_routes.route_index, block_static_routes.type, block_route_value_bools.value").
		Joins("LEFT JOIN block_static_routes on block_route_value_bools.block_route = block_static_routes.id").
		Scan(&blockRoutesBools)
	for _, route := range blockRoutesBools {
		blockRoutes = append(blockRoutes, &model.ProtoBlockStaticRoute{
			BlockStaticRoute: route.BlockStaticRoute,
			Value:            route.Value,
		})
	}
	return blockRoutes, nil
}
