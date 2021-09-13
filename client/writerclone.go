package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"strconv"
)

// ClientGetWriterClone an object
func (a *FlowClient) ClientGetWriterClone(uuid string) (*model.WriterClone, error) {
	resp, err := a.client.R().
		SetResult(&model.WriterClone{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/writers/clone/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.WriterClone), nil
}

// ClientEditWriterClone edit an object
func (a *FlowClient) ClientEditWriterClone(uuid string, body model.WriterClone, updateProducer bool) (*model.WriterClone, error) {
	param := strconv.FormatBool(updateProducer)
	resp, err := a.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": uuid}).
		SetQueryParam("update_producer", param).
		Patch("/api/writers/clone/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.WriterClone), nil
}