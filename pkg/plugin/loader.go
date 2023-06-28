package plugin

import (
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/config"
)

// / plugin loader, you need build with -buildmode=plugin and make sure there are onInit function in plugin that will be called when plugin loaded
var Plugins = make(map[string]*plugin.Plugin)

func Load() {
	pluginPath := config.C.PluginPath
	// load plugin
	pluginDirs, err := os.ReadDir(pluginPath)
	if err != nil {
		logrus.Fatalf("read plugin path %s error: %v", pluginPath, err)
		return
	}
	for _, file := range pluginDirs {
		if !file.IsDir() {
			continue
		}
		pluginDir, err := os.ReadDir(path.Join(pluginPath, file.Name()))
		if err != nil {
			logrus.Errorf("read plugin path %s error: %v", path.Join(pluginPath, file.Name()), err)
			continue
		}
		for _, pluginFile := range pluginDir {
			if !strings.HasSuffix(pluginFile.Name(), ".so") {
				continue
			}
			LoadPlugin(path.Join(pluginPath, file.Name(), pluginFile.Name()))
		}
	}
}

func LoadPlugin(pluginPath string) {
	logrus.Infof("load plugin %s", pluginPath)
	pluginInstance, err := plugin.Open(pluginPath)
	if err != nil {
		logrus.Errorf("load plugin %s error: %v", pluginPath, err)
		return
	}
	Plugins[pluginPath] = pluginInstance
	// load plugin
	symbol, err := pluginInstance.Lookup("OnInit")
	if err != nil {
		logrus.Errorf("load plugin %s error: %v", pluginPath, err)
		return
	}
	// call onLoad function
	onLoad, ok := symbol.(func())
	if !ok {
		logrus.Errorf("load plugin %s error: onLoad function not found", pluginPath)
		return
	}
	onLoad()
	logrus.Infof("load plugin %s success", pluginPath)
}

func Stop() {
	for pluginPath, pluginInstance := range Plugins {
		// load plugin
		symbol, err := pluginInstance.Lookup("OnStop")
		if err != nil {
			logrus.Errorf("load plugin %s error: %v", pluginPath, err)
			return
		}
		// call onLoad function
		onStop, ok := symbol.(func())
		if !ok {
			logrus.Errorf("load plugin %s error: onStop function not found", pluginPath)
			return
		}
		onStop()
		logrus.Infof("load plugin %s success", pluginPath)
	}
}
