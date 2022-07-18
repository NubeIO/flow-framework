package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetPoints(args api.Args) ([]*model.Point, error) {
	q, err := getDb().GetPoints(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPoint(uuid string, args api.Args) (*model.Point, error) {
	q, err := getDb().GetPoint(uuid, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreatePoint(body *model.Point, fromPlugin, updatePoint bool) (*model.Point, error) {
	pnt, err := getDb().CreatePoint(body, fromPlugin)
	if err != nil {
		return nil, err
	}
	if updatePoint {
		pnt, err = getDb().UpdatePoint(pnt.UUID, pnt, fromPlugin) // MARC: UpdatePoint is called here so that the PresentValue and Priority are updated to use the fallback value.  Otherwise they are left as Null and the Edge28 Outputs are left floating.
		if err != nil {
			return nil, err
		}
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

func (h *Handler) WritePoint(uuid string, body *model.PointWriter, fromPlugin bool) (returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	q, isPresentValueChange, isWriteValueChange, isPriorityChanged, err := getDb().PointWrite(uuid, body, fromPlugin)
	if err != nil {
		return nil, false, false, false, err
	}
	return q, isPresentValueChange, isWriteValueChange, isPriorityChanged, nil
}

func (h *Handler) UpdatePointValue(uuid string, body *model.Point, priority *map[string]*float64, fromPlugin bool) (returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	var pointModel *model.Point
	query := getDb().DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, false, false, false, query.Error
	}
	_ = getDb().DB.Model(&pointModel).Updates(&body)
	p, isPresentValueChange, isWriteValueChange, isPriorityChanged, err := getDb().UpdatePointValue(pointModel, priority, fromPlugin)
	if err != nil {
		return nil, false, false, false, err
	}

	return p, isPresentValueChange, isWriteValueChange, isPriorityChanged, nil
}

func (h *Handler) GetOnePointByArgs(args api.Args) (*model.Point, error) {
	return getDb().GetOnePointByArgs(args)
}

func (h *Handler) DeletePoint(uuid string) (bool, error) {
	_, err := getDb().DeletePoint(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
