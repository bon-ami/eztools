package eztools

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
)

// XMLRead reads config file into input structure
func XMLRead(name string, data interface{}) error {
	// directly open, without checking existence
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	/*decoder := xml.NewDecoder(file)
	err = decoder.Decode(data)
	ShowSthln(data)*/
	bytes, _ := ioutil.ReadAll(file)
	err = xml.Unmarshal(bytes, data)
	return err
}

// XMLsReadDefault reads config file into input structure from
// give path, or give file name (plus .xml) under current dir or home dir
func XMLsReadDefault(path, file string, cfg interface{}) (err error) {
	if len(path) > 0 {
		err = XMLRead(path, cfg)
		if err == nil {
			return
		}
	}
	if len(file) < 1 {
		return ErrNoValidResults
	}
	home, _ := os.UserHomeDir()
	cfgPaths := [...]string{".", home}
	for _, path1 := range cfgPaths {
		err = XMLRead(filepath.Join(path1, file+".xml"), cfg)
		if err == nil {
			break
		}
	}
	return
}
