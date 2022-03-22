package dbhandler

import (
	"fmt"
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

func (h *Handler) GetPoint(uuid string, args api.Args) (*model.Point, error) {
	q, err := getDb().GetPoint(uuid, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreatePoint(body *model.Point, fromPlugin, updatePoint bool) (*model.Point, error) {
	fmt.Println("CreatePoint() BEFORE")
	fmt.Println("CreatePoint() BEFORE: point")
	fmt.Printf("%+v\n", body)
	pnt, err := getDb().CreatePoint(body, fromPlugin)
	if err != nil {
		fmt.Println("CreatePoint() ERROR 1")
		return nil, err
	}
	fmt.Println("CreatePoint() AFTER: pnt")
	fmt.Printf("%+v\n", pnt)
	fmt.Println("CreatePoint() updatePoint: ", updatePoint, " fromPlugin: ", fromPlugin)
	if updatePoint {
		fmt.Println("UpdatePoint() in CreatePoint() BEFORE")
		pnt, err = getDb().UpdatePoint(pnt.UUID, pnt, fromPlugin) //MARC: UpdatePoint is called here so that the PresentValue and Priority are updated to use the fallback value.  Otherwise they are left as Null and the Edge28 Outputs are left floating.
		if err != nil {
			fmt.Println("UpdatePoint() in CreatePoint() ERROR")
			return nil, err
		}
		fmt.Println("UpdatePoint() in CreatePoint() AFTER")
		fmt.Println("UpdatePoint() in CreatePoint() AFTER: pnt")
		fmt.Printf("%+v\n", pnt)
	}
	return pnt, nil
}

func (h *Handler) UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	fmt.Println("DB HANDLER UpdatePoint(): body")
	fmt.Printf("%+v\n", body)
	q, err := getDb().UpdatePoint(uuid, body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdatePointPresentValue(body *model.Point, fromPlugin bool) (*model.Point, error) {
	p, err := getDb().UpdatePointValue(body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (h *Handler) UpdatePointValue(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := getDb().DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	_ = getDb().DB.Model(&pointModel).Updates(&body)
	// Don't update point value if priority array on body is nil
	if body.Priority == nil {
		return pointModel, nil
	} else {
		pointModel.Priority = body.Priority
	}
	pointModel.Priority = body.Priority
	p, err := getDb().UpdatePointValue(pointModel, fromPlugin)
	if err != nil {
		return nil, err
	}

	return p, nil
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
