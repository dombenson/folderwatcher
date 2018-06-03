package iniconfig

import (
	"dgeb/config"
	"dgeb/defaults"
	"dgeb/interfaces"
	"log"
	"os"
	"os/user"

	"github.com/dombenson/go-ini"
)

// Get configuration from ini file
func Get(preferredFile string) interfaces.Config {
	c := config.Create()

	// Find config file options
	possibleConfigFiles := findCandidateFilePaths(preferredFile)

	// Load in turn: so we have applied internal defaults,
	// then overlay global config, then overlay user config,
	// finally overlay any params in an explicitly provided config
	for _, filePath := range possibleConfigFiles {
		processOneFileCandidate(c, filePath)
	}

	return c
}

func findCandidateFilePaths(preferredFile string) []string {
	var userHomeDir string
	curUser, err := user.Current()
	if err != nil {
		log.Println("Unable to get current user for config search")
	} else {
		userHomeDir = curUser.HomeDir
	}

	possibleConfigFiles := make([]string, 0, 3)
	possibleConfigFiles = append(possibleConfigFiles, defaults.GlobalConfigFile)
	if userHomeDir != "" {
		possibleConfigFiles = append(possibleConfigFiles, userHomeDir+"/"+defaults.UserConfigFile)
	}
	if preferredFile != "" {
		possibleConfigFiles = append(possibleConfigFiles, preferredFile)
	}
	return possibleConfigFiles
}

func processOneFileCandidate(c *config.Config, filePath string) {
	fh, err := os.Open(filePath)
	if os.IsNotExist(err) {
		// Config candidate not existing is not a reportable error
		return
	} else if err != nil {
		log.Println("Open config file error:", filePath, err)
		return
	}
	defer fh.Close()
	thisIni, err := ini.Load(fh)
	if err != nil {
		log.Println("Error loading data from ini file", filePath, err)
	}
	defer thisIni.Close()

	log.Printf("Adding config from '%s'\n", filePath)
	applyConfigFromIni(c, thisIni)

}

func applyConfigFromIni(c *config.Config, thisIni ini.File) {
	var newStrVal string
	var newIntVal int
	var ok bool
	newStrVal, ok = thisIni.Get("client", "discover_addr")
	if ok {
		c.ClientDiscover = newStrVal
	}
	newIntVal, ok = thisIni.GetInt("client", "discover_port")
	if ok {
		c.ClientDiscoverPort = newIntVal
	}
	newStrVal, ok = thisIni.Get("client", "path")
	if ok {
		c.MonitorPath = newStrVal
	}

	newStrVal, ok = thisIni.Get("server", "discover_addr")
	if ok {
		c.ServerDiscover = newStrVal
	}
	newIntVal, ok = thisIni.GetInt("server", "discover_port")
	if ok {
		c.ServerDiscoverPort = newIntVal
	}

	newIntVal, ok = thisIni.GetInt("server", "http_port")
	if ok {
		c.HTTPPort = newIntVal
	}
}
