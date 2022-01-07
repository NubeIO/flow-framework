package jobs

import log "github.com/sirupsen/logrus"

func (j *Jobs) RefreshTokenJobAdd() error {
	_, err := cron.Every(1).Hour().Tag("refreshToken").Do(j.refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jobs) refreshToken() {
	log.Info("REFRESH TOKEN RUN")
	_, err := j.db.RefreshLocalStorageFlowToken()
	if err != nil {
		log.Error(err)
	}
	_, err = j.db.RefreshFlowNetworksConnections()
	if err != nil {
		log.Error(err)
	}
	_, err = j.db.RefreshFlowNetworkClonesConnections()
	if err != nil {
		log.Error(err)
	}
}
