package pkg

import (
	"build-booster/common"
	"build-booster/common/blog"
	"build-booster/common/http/httpserver"
	"build-booster/server/pkg/resource/direct/agent/config"
	"build-booster/server/pkg/resource/direct/agent/pkg/api"
	"build-booster/server/pkg/resource/direct/agent/pkg/types"
)

// FbAgent : fast build agent
type FbAgent struct {
	conf       *config.ServerConfig
	httpServer *httpserver.HTTPServer
	handle     *api.HTTPHandle
}

// NewFbAgent : return fast build agent object
func NewFbAgent(conf *config.ServerConfig) (*FbAgent, error) {
	s := &FbAgent{conf: conf}

	// Http server
	s.httpServer = httpserver.NewHTTPServer(s.conf.Port, s.conf.Address, "")
	if s.conf.ServerCert.IsSSL {
		s.httpServer.SetSSL(
			s.conf.ServerCert.CAFile, s.conf.ServerCert.CertFile, s.conf.ServerCert.KeyFile, s.conf.ServerCert.CertPwd)
	}

	s.initConfig()
	return s, nil
}

// init same inner config by ServerConfig
func (server *FbAgent) initConfig() {
	types.DistCCDaemonCPUPerUnit = server.conf.BcsCPUPerInstance
}

// Start : start listen
func (server *FbAgent) Start() error {
	var err error
	server.handle, err = api.NewHTTPHandle(server.conf)
	if server.handle == nil || err != nil {
		return types.ErrInitHTTPHandle
	}

	server.httpServer.RegisterWebServer(api.PathV1, nil, server.handle.GetActions())

	return server.httpServer.ListenAndServe()
}

// Run brings up the server
func Run(conf *config.ServerConfig) error {
	if err := common.SavePid(conf.ProcessConfig); err != nil {
		blog.Errorf("save pid failed: %v", err)
		return err
	}

	server, err := NewFbAgent(conf)
	if err != nil {
		blog.Errorf("init distCC server failed: %v", err)
		return err
	}

	return server.Start()
}
