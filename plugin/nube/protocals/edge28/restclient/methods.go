package edgerest

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ServerPing struct {
	State string `json:"1_state"`
}

type WriteResponseUO struct {
	State string  `json:"1_state"`
	IoNum string  `json:"2_ioNum"`
	Gpio  string  `json:"3_gpio"`
	Val   float64 `json:"4_val"`
	Msg   string  `json:"5_msg"`
}

type WriteResponse struct {
	State string `json:"1_state"`
	IoNum string `json:"2_ioNum"`
	Gpio  string `json:"3_gpio"`
	Val   int    `json:"4_val"`
	Msg   string `json:"5_msg"`
}

type UI struct {
	State string `json:"1_state"`
	IoNum string `json:"2_ioNum"`
	Gpio  string `json:"3_gpio"`
	Val   struct {
		UI1 struct {
			Val float64 `json:"val"`
		} `json:"UI1"`
		UI2 struct {
			Val float64 `json:"val"`
		} `json:"UI2"`
		UI3 struct {
			Val float64 `json:"val"`
		} `json:"UI3"`
		UI4 struct {
			Val float64 `json:"val"`
		} `json:"UI4"`
		UI5 struct {
			Val float64 `json:"val"`
		} `json:"UI5"`
		UI6 struct {
			Val float64 `json:"val"`
		} `json:"UI6"`
		UI7 struct {
			Val float64 `json:"val"`
		} `json:"UI7"`
	} `json:"4_val"`
	Msg      string `json:"5_msg"`
	MinRange struct {
		UI1 int `json:"UI1"`
		UI2 int `json:"UI2"`
		UI3 int `json:"UI3"`
		UI4 int `json:"UI4"`
		UI5 int `json:"UI5"`
		UI6 int `json:"UI6"`
		UI7 int `json:"UI7"`
	} `json:"6_min_range"`
	MaxRange struct {
		UI1 int `json:"UI1"`
		UI2 int `json:"UI2"`
		UI3 int `json:"UI3"`
		UI4 int `json:"UI4"`
		UI5 int `json:"UI5"`
		UI6 int `json:"UI6"`
		UI7 int `json:"UI7"`
	} `json:"7_max_range"`
}

type DI struct {
	State string `json:"1_state"`
	IoNum string `json:"2_ioNum"`
	Gpio  string `json:"3_gpio"`
	Val   struct {
		DI1 struct {
			Val float64 `json:"val"`
		} `json:"DI1"`
		DI2 struct {
			Val float64 `json:"val"`
		} `json:"DI2"`
		DI3 struct {
			Val float64 `json:"val"`
		} `json:"DI3"`
		DI4 struct {
			Val float64 `json:"val"`
		} `json:"DI4"`
		DI5 struct {
			Val float64 `json:"val"`
		} `json:"DI5"`
		DI6 struct {
			Val float64 `json:"val"`
		} `json:"DI6"`
		DI7 struct {
			Val float64 `json:"val"`
		} `json:"DI7"`
	} `json:"4_val"`
	Msg string `json:"5_msg"`
}

// PingServer all points
func (a *RestClient) PingServer() (*ServerPing, error) {
	resp, err := a.client.R().
		SetResult(&ServerPing{}).
		Get("/")
	if err != nil {
		return nil, fmt.Errorf("error geting server %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ServerPing), nil
}

// GetUIs all ui points
func (a *RestClient) GetUIs() (*UI, error) {
	resp, err := a.client.R().
		SetResult(UI{}).
		Get("/api/1.1/read/all/ui")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*UI), nil
}

// GetDIs all di points
func (a *RestClient) GetDIs() (*DI, error) {
	resp, err := a.client.R().
		SetResult(DI{}).
		Get("/api/1.1/read/all/di")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*DI), nil
}

// WriteUO to a UO point
func (a *RestClient) WriteUO(ioNum string, value float64) (*WriteResponseUO, error) {
	req := fmt.Sprintf("/api/1.1/write/uo/%s/%d/16", strings.ToLower(ioNum), utils.ToInt(value))
	log.Println("edge28-api: WriteUO:", req)
	resp, err := a.client.R().
		SetResult(WriteResponseUO{}).
		Get(req)
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*WriteResponseUO), nil
}

// WriteDO to a DO point
func (a *RestClient) WriteDO(ioNum string, value float64) (*WriteResponse, error) {
	req := fmt.Sprintf("/api/1.1/write/do/%s/%d/16", strings.ToLower(ioNum), utils.ToInt(value))
	log.Println("edge28-api: WriteDO:", req)
	resp, err := a.client.R().
		SetResult(WriteResponse{}).
		Get(req)
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*WriteResponse), nil
}
