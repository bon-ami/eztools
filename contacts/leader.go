package contacts

import (
	"database/sql"
	"strconv"

	"github.com/bon-ami/eztools"
)

/* used by leaders only */
// return values: [][0] ID, [][1] Name
func GetMembers(db *sql.DB, id int) ([][]string, error) {
	if id == eztools.InvalidID {
		return nil, eztools.ErrInvalidInput
	}
	initFieldNames()
	team, err := eztools.Search(db, eztools.TblCONTACTS,
		eztools.FldID+"="+strconv.Itoa(id), []string{fieldNames[2]}, "")
	if err != nil {
		return nil, err
	}
	if len(team) < 1 {
		return nil, eztools.ErrNoValidResults
	}
	leader, err := eztools.Search(db, eztools.TblTEAM,
		eztools.FldID+"="+team[0][0],
		[]string{eztools.FldLEADER}, "")
	if err != nil {
		return nil, err
	}
	if len(leader) < 1 || leader[0][0] != strconv.Itoa(id) {
		return nil, eztools.ErrInvalidInput
	}
	members, err := eztools.Search(db, eztools.TblCONTACTS,
		fieldNames[2]+"="+team[0][0],
		[]string{eztools.FldID, eztools.FldNAME}, "")
	return members, err
}

func MakeLeader(db *sql.DB, id, team int) error {
	var (
		teamMembers, teamID [][]string
		err                 error
	)
	initFieldNames()
	if id < 0 || id != eztools.InvalidID {
		teamMembers, err = eztools.Search(db, eztools.TblCONTACTS,
			fieldNames[2]+"="+strconv.Itoa(team), []string{eztools.FldID, "name"}, "")
		if err != nil {
			return err
		}
	}
	teamID, err = eztools.Search(db, eztools.TblTEAM,
		eztools.FldID+"="+strconv.Itoa(team), []string{eztools.FldID, "leader"}, "")
	if err != nil {
		return err
	}
	lt := len(teamID)
	switch {
	case lt <= 0:
		eztools.ShowStrln("No team found!")
		return nil
	case lt > 1:
		eztools.ShowStrln("Multiple teams found with same ID!")
		return nil
	}
	newID := id
	if id < 0 || id != eztools.InvalidID {
		eztools.ShowStr("Members of team: ")
		for _, member := range teamMembers {
			if teamID[0][1] == member[0] {
				eztools.ShowStr("[leader]")
			}
			eztools.ShowStr(member[0] + " " + member[1] + ", ")
		}
		newID, err = eztools.PromptInt("# of new leader=")
		if err != nil {
			newID = eztools.InvalidID
		}
	}
	if newID >= 0 && newID != eztools.InvalidID {
		teamMembers, err = eztools.Search(db, eztools.TblCONTACTS,
			eztools.FldID+"="+strconv.Itoa(newID), nil, "")
		if err != nil || len(teamMembers) < 1 {
			eztools.ShowStrln("invalid ID")
			return nil
		}
		err = eztools.UpdateWtParams(db, eztools.TblTEAM,
			eztools.FldID+"="+teamID[0][0],
			[]string{eztools.FldLEADER},
			[]string{strconv.Itoa(newID)},
			false)
	}
	return err
}
