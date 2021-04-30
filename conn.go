package eztools

import (
	"database/sql"
	"encoding/xml"
)

const (
	strSQL     = "mysql"
	strCfgFile = "eztools"
)

// ConnectWtParam connects to the database using parameters.
func ConnectWtParam(user, pass, ip, database string) (db *sql.DB, err error) {

	db, err = sql.Open(strSQL, user+":"+pass+"@"+ip+"/"+database)
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return
}

// ConnectWtCfg connects to the database using parameters from an xml file
// root element is named "root", elements include "ip", "db", "user" and "pass"
func ConnectWtCfg(file string) (*sql.DB, error) {
	var cfg struct {
		Root    xml.Name `xml:"root"`
		StrUSER string   `xml:"user"`
		StrPASS string   `xml:"pass"`
		StrIP   string   `xml:"ip"`
		StrDB   string   `xml:"db"`
	}

	if _, err := XMLsReadDefaultNoCreate("", file, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.StrUSER) < 1 ||
		len(cfg.StrDB) < 1 ||
		len(cfg.StrIP) < 1 ||
		len(cfg.StrPASS) < 1 {
		return nil, ErrInvalidInput
	}
	return ConnectWtParam(cfg.StrUSER, cfg.StrPASS, cfg.StrIP, cfg.StrDB)
}

// Connect connects to the database using parameters from eztools.xml
// root element is named "root", elements include "ip", "db", "user" and "pass"
func Connect() (*sql.DB, error) {
	return ConnectWtCfg(strCfgFile)
}
