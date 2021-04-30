package eztools

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
)

// XMLWrite writes input structure to file
func XMLWrite(file string, data interface{}, createIfNeeded bool) error {
	if !createIfNeeded {
		if _, err := os.Stat(file); err != nil {
			return err
		}
	}
	bytes, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, bytes, 0664)
}

// XMLWriteNoCreate writes config file from input structure
// by a full path (with .xml extension)
func XMLWriteNoCreate(file string, cfg interface{}) error {
	return XMLWrite(file, cfg, false)
}

// XMLRead reads file into input structure
func XMLRead(file string, data interface{}, _ bool) error {
	if _, err := os.Stat(file); err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return xml.Unmarshal(bytes, data)
}

// XMLsReadDefault reads a config file
func XMLsReadDefault(path, file string, cfg interface{},
	createIfNeeded bool) (
	pathFound string, err error) {
	if len(path) > 0 {
		err = XMLRead(path, cfg, createIfNeeded)
		if err == nil {
			return path, err
		}
	}
	if len(file) < 1 {
		return "", ErrNoValidResults
	}
	home, _ := os.UserHomeDir()
	cfgPaths := [...]string{".", home}
	for _, path1 := range cfgPaths {
		pathFound = filepath.Join(path1, file+".xml")
		err = XMLRead(pathFound, cfg, createIfNeeded)
		if err == nil {
			return
		}
	}
	return "", ErrNoValidResults
}

// XMLsReadDefaultNoCreate reads config file into input structure from
// given path, or given file name (plus .xml) under current dir or home dir
// returns full file name with path
func XMLsReadDefaultNoCreate(path, file string, cfg interface{}) (string, error) {
	return XMLsReadDefault(path, file, cfg, false)
}
