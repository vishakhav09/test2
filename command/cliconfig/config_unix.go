// +build !windows

package cliconfig

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func configFile() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}

	// If legacy ~/.terraformrc dir exists already, prefer that
	file := filepath.Join(home, ".terraformrc")
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		log.Printf("[DEBUG] Found .terraformrc in legacy location, continuing")
		return file, nil
	}

	// else use configDir()'s result and attach terraformrc
	dir, err := configDir()
	if err != nil {
		return "", err
	}

	// With the legacy dir still keep the dot infront of terraformrc
	legacyDir, err := legacyConfigDir()
	if dir == legacyDir {
		return filepath.Join(dir, ".terraformrc"), nil
	}

	return filepath.Join(dir, "terraformrc"), nil
}

func configDir() (string, error) {
	log.Printf("[DEBUG] Getting config dir")

	// First prefer the XDG_CONFIG_HOME environmental variable
	configDirPath := os.Getenv("XDG_CONFIG_HOME")
	if configDirPath != "" {
		log.Printf("[DEBUG] Prefering XDG_CONFIG_HOME env variable")
		return filepath.Join(configDirPath, "terraform"), nil
	}

	// If legacy ~/.terraform.d dir exists already, prefer that
	configDirPath, err := legacyConfigDir()
	if err == nil {
		if _, err := os.Stat(configDirPath); !os.IsNotExist(err) {
			log.Printf("[DEBUG] Found .terraform.d directory in legacy location, continuing")
			return configDirPath, nil
		}
	}
	// Else fall back to XDG_CONFIG_HOME's standard location $HOME/.config
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Falling back to XDG_CONFIG_HOME standard location $HOME/.config")
	return filepath.Join(home, ".config", "terraform"), nil

}

func legacyConfigDir() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".terraform.d"), nil
}

func cacheDir() (string, error) {
	log.Printf("[DEBUG] Getting cache dir")

	// First prefer the XDG_CACHE_HOME environmental variable
	cacheDirPath := os.Getenv("XDG_CACHE_HOME")
	if cacheDirPath != "" {
		log.Printf("[DEBUG] Prefering XDG_CACHE_HOME env variable")
		return filepath.Join(cacheDirPath, "terraform"), nil
	}

	// If legacy ~/.terraform.d dir exists already, prefer that
	cacheDirPath, err := legacyConfigDir()
	if err == nil {
		if _, err := os.Stat(cacheDirPath); !os.IsNotExist(err) {
			log.Printf("[DEBUG] Found .terraform.d directory in legacy location, continuing")
			return cacheDirPath, nil
		}
	}
	// Else fall back to XDG_CACHE_HOME's standard location $HOME/.cache
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Falling back to XDG_CACHE_HOME standard location $HOME/.cache")
	return filepath.Join(home, ".cache", "terraform"), nil
}

func homeDir() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		// FIXME: homeDir gets called from globalPluginDirs during init, before
		// the logging is setup.  We should move meta initializtion outside of
		// init, but in the meantime we just need to silence this output.
		//log.Printf("[DEBUG] Detected home directory from env var: %s", home)

		return home, nil
	}

	// If that fails, try build-in module
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	if user.HomeDir == "" {
		return "", errors.New("blank output")
	}

	return user.HomeDir, nil
}

func replaceFileAtomic(source, destination string) error {
	// On Unix systems, a rename is sufficiently atomic.
	return os.Rename(source, destination)
}
