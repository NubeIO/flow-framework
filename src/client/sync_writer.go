package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (a *FlowClient) SyncWriter(body *model.SyncWriter) (*model.WriterClone, error) {
	resp, err := a.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		Post("/api/sync/writer")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("SyncWriter: %s", err)
		} else {
			return nil, fmt.Errorf("SyncWriter: %s", resp)
		}
	}
	return resp.Result().(*model.WriterClone), nil
}

func (a *FlowClient) SyncCOV(writerUUID string, body *model.SyncCOV) error {
	resp, err := a.client.R().
		SetBody(body).
		Post(fmt.Sprintf("/api/sync/csv/%s", writerUUID))
	if err != nil {
		if resp == nil || resp.String() == "" {
			return fmt.Errorf("SyncCOV: %s", err)
		} else {
			return fmt.Errorf("SyncCOV: %s", resp)
		}
	}
	return nil
}

func (a *FlowClient) SyncWriterWriteAction(sourceUUID string, body *model.SyncWriterAction) error {
	resp, err := a.client.R().
		SetBody(body).
		Post(fmt.Sprintf("/api/sync/writer/write/%s", sourceUUID))
	// TODO: this block needs to be re-written; same constant thing on all places
	if err != nil {
		if resp == nil || resp.String() == "" {
			return fmt.Errorf("SyncWriterWriteAction: %s", err)
		} else {
			return fmt.Errorf("SyncWriterWriteAction: %s", resp)
		}
	} else if !(resp.StatusCode() >= 200 && resp.StatusCode() < 300) {
		return fmt.Errorf("SyncWriterWriteAction: %s", resp)
	}
	return nil
}

func (a *FlowClient) SyncWriterReadAction(sourceUUID string) error {
	resp, err := a.client.R().Get(fmt.Sprintf("/api/sync/writer/read/%s", sourceUUID))
	if err != nil {
		if resp == nil || resp.String() == "" {
			return fmt.Errorf("SyncWriterReadAction: %s", err)
		} else {
			return fmt.Errorf("SyncWriterReadAction: %s", resp)
		}
	}
	return nil
}
