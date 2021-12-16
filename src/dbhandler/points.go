package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
)

func (h *Handler) GetPoints(args api.Args) ([]*model.Point, error) {
	q, err := getDb().GetPoints(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPoint(uuid string) (*model.Point, error) {
	q, err := getDb().GetPoint(uuid, api.Args{})
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreatePoint(body *model.Point) (*model.Point, error) {
	pnt, err := getDb().CreatePoint(body, "")
	if err != nil {
		return nil, err
	}
	pnt, err = getDb().UpdatePoint(pnt.UUID, pnt, false) //MARC: UpdatePoint is called here so that the PresentValue and Priority are updated to use the fallback value.  Otherwise they are left as Null and the Edge28 Outputs are left floating.
	if err != nil {
		return nil, err
	}
	return pnt, nil
}

func (h *Handler) UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	q, err := getDb().UpdatePoint(uuid, body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdatePointValue(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	q, err := getDb().UpdatePointValue(uuid, body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPointByField(field string, value string) (*model.Point, error) {
	q, err := getDb().GetPointByField(field, value)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdatePointByFieldAndUnit(field string, value string, body *model.Point) (*model.Point, error) {
	q, err := getDb().UpdatePointByFieldAndUnit(field, value, body, false)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPointByFieldAndIOID(field string, value string, body *model.Point) (*model.Point, error) {
	q, err := getDb().GetPointByFieldAndIOID(field, value, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPointByFieldAndThingType(field string, value string, body *model.Point) (*model.Point, error) {
	q, err := getDb().GetPointByFieldAndThingType(field, value, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) PointAndQuery(value1, value2 string) (*model.Point, error) {
	q, err := getDb().PointAndQuery(value1, value2)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) DeletePoint(uuid string) (bool, error) {
	_, err := getDb().DeletePoint(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
