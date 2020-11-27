package teams

import (
	"database/sql"

	"github.com/bon-ami/eztools"
)

func List(db *sql.DB, team int) error {
	teamArr, err := eztools.Search(db, eztools.TblTEAM, "",
		[]string{eztools.FldID, eztools.FldSTR, eztools.FldLEADER}, "")
	if err != nil {
		return err
	}
	var (
		contactArr [][]string
		leader     string
	)
	for _, team := range teamArr {
		contactArr, err = eztools.Search(db, eztools.TblCONTACTS,
			eztools.FldID+"="+team[2], []string{eztools.FldNAME}, "")
		if err != nil {
			leader = "None"
		} else {
			if len(contactArr) == 0 {
				leader = "Invalid"
			} else {
				leader = contactArr[0][0]
			}
		}
		eztools.ShowStrln(team[0] + ": " + team[1] + " leader: " + leader)
	}
	return nil
}
