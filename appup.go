package eztools

import (
	"database/sql"

	"github.com/sanbornm/go-selfupdate/selfupdate"
)

// AppUpgrade checks for updates and applies the update automatically,
// which will work next time the app is run.
func AppUpgrade(db *sql.DB, prefix string, ver string, server *chan string, ch chan bool) {
	type pairs struct {
		id  string
		str string
	}
	upPairs := []pairs{
		pairs{
			"Url", "",
		},
		pairs{
			"Dir", "",
		},
		pairs{
			"App", "",
		},
	}
	if ver == "dev" {
		upPairs[0].id += "Dev"
	}
	var (
		upStr string
		err   error
	)
	for i, upPair := range upPairs {
		upStr, err = GetPairStr(db, TblCHORE, prefix+upPair.id)
		if err == nil && len(upStr) > 0 {
			upPairs[i].str = upStr
		} else {
			LogPrint("NO " + prefix + upPair.id + " configured!")
			if err != nil {
				LogErrPrint(err)
			}
			break
		}
	}
	if len(upPairs[0].str) < 1 {
		Log("update check response mal-configured")
		ch <- false
		return
	}
	Log("update check response temp")
	ch <- true
	if server != nil {
		*server <- upPairs[0].str
	}
	if len(upPairs[1].str) < 1 || len(upPairs[2].str) < 1 {
		ch <- false
		Log("update check response final failure")
		return
	}
	var updater = &selfupdate.Updater{
		CurrentVersion: ver,
		ApiURL:         upPairs[0].str,
		BinURL:         upPairs[0].str,
		DiffURL:        upPairs[0].str,
		Dir:            upPairs[1].str,
		CmdName:        upPairs[2].str,
		ForceCheck:     true,
	}
	if ver == "ro" || len(ver) < 1 {
		ShowStrln("update skipped for ro or nil versions.")
	} else {
		if err = updater.BackgroundRun(); err != nil {
			ShowStrln("Failed to check update for this app")
			LogErr(err)
		} else {
			if len(updater.Info.Version) > 0 && ver != updater.Info.Version {
				ShowStrln("On the next run, this app will be updated to V" + updater.Info.Version)
			}
		}
	}
	Log("update check response final")
	ch <- true
}
