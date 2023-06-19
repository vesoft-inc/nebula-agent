package license

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/unknwon/goconfig"
	"github.com/vesoft-inc/license/pkg/lmclient"
	pkgtypes "github.com/vesoft-inc/license/pkg/types"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
)

const (
	LicenseManagerSDConfigPath = "download/prometheus/sd_config/license_manager.yml"
)

var LicenseApplyHandler = new(licenseApplyHandler)

type FileSDConfigItem struct {
	Targets []string               `yaml:"targets"`
	Labels  map[string]interface{} `yaml:"labels,omitempty"`
}

type licenseApplyHandler struct {
	Initing       bool                         `json:"init"`
	IsStopService bool                         `json:"isStopService"`
	Result        *pkgtypes.LicenseApplyResult `json:"result"`
	Error         string                       `json:"error"`
	URL           string                       `json:"url"`
}

func (h *licenseApplyHandler) GetApplyData() *pkgtypes.LicenseApplyData {
	return &pkgtypes.LicenseApplyData{
		Name: pkgtypes.LicenseProductNameDashboard,
	}
}

func (h *licenseApplyHandler) Handle(resp *lmclient.LicenseApplyHandleData) {
	h.IsStopService = resp.IsStopService
	h.Result = resp.Result
	h.Error = ""
	h.Initing = false
	if resp.Error != nil {
		h.Error = resp.Error.Error()
		logrus.Errorf("License apply handle error: %s", resp.Error)
	}
}

func (h *licenseApplyHandler) LoadEncryptedLicenseCache() (string, error) {
	ciphertext, err := os.ReadFile(".license.cache")
	return string(ciphertext), err
}

func (h *licenseApplyHandler) UpdateEncryptedLicenseCache(ciphertext string) error {
	return os.WriteFile(".license.cache", []byte(ciphertext), 0600)
}

func (h *licenseApplyHandler) Infof(format string, a ...any) {
	logrus.Infof("[License Apply]"+format, a...)
}

func (h *licenseApplyHandler) Errorf(format string, a ...any) {
	logrus.Errorf("[License Apply]"+format, a...)
}

var applyTimer lmclient.LicenseApplyTimer
var lmclientInstance lmclient.Client

func InitApplyTimer() error {
	LicenseApplyHandler.Initing = true
	LicenseApplyHandler.IsStopService = true
	LicenseApplyHandler.Error = "initing"
	lmConfigPath := path.Join(config.C.AnalyticsPath, "scripts/analytics.conf")
	config, err := goconfig.LoadConfigFile(lmConfigPath)
	if err != nil {
		LicenseApplyHandler.Error = err.Error()
		return err
	}
	logrus.Info("keys:", config.GetKeyList(""))
	lmURL, err := config.GetValue("", "--license_manager_url")
	if err != nil {
		LicenseApplyHandler.Error = err.Error()
		return err
	}
	if applyTimer != nil {
		applyTimer.Stop()
		applyTimer = nil
	}
	lmclientInstance = lmclient.New(lmURL)
	applyTimer = lmclient.NewLicenseApplyTimer(lmclientInstance, LicenseApplyHandler)
	applyTimer.Start()
	logrus.Infof("License apply timer started")
	for {
		if !LicenseApplyHandler.Initing {
			break
		}
	}
	return nil
}
