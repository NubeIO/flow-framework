package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

// FlowClient is used to invoke Form3 Accounts API.
type FlowClient struct {
	client      *resty.Client
	ClientToken string
}

func NewFlowClientCli(ip *string, port *int, token *string, isMasterSlave *bool, globalUUD string, isFNCreator bool) *FlowClient {
	if utils.IsTrue(isMasterSlave) {
		if isFNCreator {
			return newSlaveToMasterCallSession()
		} else {
			return newMasterToSlaveSession(globalUUD)
		}
	} else {
		return newSessionWithToken(*ip, *port, *token)
	}
}

func newSessionWithToken(ip string, port int, token string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d", ip, port)
	client.SetHostURL(url)
	client.SetError(&Error{})
	client.SetHeader("Authorization", token)
	return &FlowClient{client: client}
}

func newMasterToSlaveSession(globalUUID string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d/slave/%s/ff", "0.0.0.0", 1616, globalUUID)
	client.SetHostURL(url)
	client.SetError(&Error{})
	client.SetHeader("Authorization", getRubixServiceInternalToken())
	return &FlowClient{client: client}
}

func newSlaveToMasterCallSession() *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d/master/ff", "0.0.0.0", 1616)
	client.SetHostURL(url)
	client.SetError(&Error{})
	client.SetHeader("Authorization", getRubixServiceInternalToken())
	return &FlowClient{client: client}
}

func getRubixServiceInternalToken() string {
	rubixServiceDataLocation := "/data/rubix-service" //TODO: move on registry or config
	relativeAuthDataFile := "/data/internal_token.txt"
	authDataFile := path.Join(rubixServiceDataLocation, relativeAuthDataFile)
	file, err := os.Open(authDataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	internalToken, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(internalToken)
}
