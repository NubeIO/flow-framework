package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"strconv"
)

// GetWriterClone an object
func (a *FlowClient) GetWriterClone(uuid string) (*model.WriterClone, error) {
	resp, err := a.client.R().
		SetResult(&model.WriterClone{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/producers/writer_clones/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.WriterClone), nil
}

// EditWriterClone edit an object
func (a *FlowClient) EditWriterClone(uuid string, body model.WriterClone, updateProducer bool) (*model.WriterClone, error) {
	param := strconv.FormatBool(updateProducer)
	resp, err := a.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": uuid}).
		SetQueryParam("update_producer", param).
		Patch("/api/producers/writer_clones/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.WriterClone), nil
}

// CreateWriterClone edit an object
func (a *FlowClient) CreateWriterClone(body model.WriterClone) (*model.WriterClone, error) {
	resp, err := a.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		Post("/api/producers/writer_clones")
	if err != nil {
		return nil, fmt.Errorf("%s failed", err)
	}
	fmt.Println(resp.String())
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.WriterClone), nil
}
