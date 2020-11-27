package contacts

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/bon-ami/eztools"
)

var (
	fieldNames   []string
	fieldNameStr string
)

func initFieldNames() {
	if len(fieldNames) == 0 {
		// is there a better way to initialize fieldNames
		// without fieldNamesTemplate
		fieldNamesTemplate := [...]string{
			eztools.FldNUMBER, eztools.FldNAME, eztools.FldTEAM,
			eztools.FldEXT, eztools.FldPHONE, eztools.FldMAIL,
			eztools.FldLDAP, eztools.FldUID, eztools.FldNICK}
		fieldNames = fieldNamesTemplate[:]
		fieldNameStr = strings.Join(fieldNames, ", ")
	}
}

func list1Field(i int, val string) {
	eztools.ShowStr(fieldNames[i] + "=" + val + ", ")
}

// ListFieldsFrmMap prints contact info out of the map.
//   If changedOnly is true, db and id will be used to get current values and
// prints changed info only.
func ListFieldsFrmMap(db *sql.DB, id int, flds map[string]string, changedOnly bool) (changedFlds []string) {
	var (
		searched [][]string
		err      error
	)
	idStr := strconv.Itoa(id)
	if changedOnly {
		if db == nil || id == eztools.InvalidID {
			eztools.Log("NO database or ID provided to list contact info!")
			return
		}
		searched, err = eztools.Search(db, eztools.TblCONTACTS,
			eztools.FldID+"="+idStr,
			fieldNames, "")
		if err != nil && err != eztools.ErrNoValidResults {
			eztools.LogErr(err)
			return
		}
	}
	if len(searched) > 1 {
		eztools.Log(strconv.Itoa(len(searched)) + " records found for ID " + idStr)
	}
	eztools.ShowStr("ID=" + idStr)
	for i, v := range fieldNames {
		f1 := flds[v]
		if changedOnly {
			changed := false
			for _, s1 := range searched {
				if f1 != s1[i] {
					eztools.ShowStr("\t" + v + ": " + s1[i] + " -> " + f1)
					changed = true
				}
			}
			if changed {
				changedFlds = append(changedFlds, v)
			}
		} else if len(f1) > 0 {
			eztools.ShowStr("\t" + v + ": " + f1)
		}
	}
	return
}

// ChooseField asks user to choose a field in contacts table
func ChooseField(info string) string {
	if len(info) > 0 {
		eztools.ShowStrln(info)
	}
	if ret := eztools.ChooseStrings(fieldNames); ret != eztools.InvalidID {
		return fieldNames[ret]
	}
	return ""
}

// List shows all members' info of a team
func List(db *sql.DB, team int) error {
	initFieldNames()
	var cri string
	if team > 0 {
		cri = fieldNames[2] + "=" + strconv.Itoa(team)
	}
	selected, err := eztools.Search(db, eztools.TblCONTACTS, cri,
		append(fieldNames, eztools.FldID), "")
	if err != nil {
		return err
	}
	teamArr, err := eztools.Search(db, eztools.TblTEAM, "",
		[]string{eztools.FldID, eztools.FldSTR, eztools.FldLEADER}, "")
	if err != nil {
		return err
	}
	type teamCont struct {
		str, leader string
	}
	teamMap := make(map[string]teamCont)
	for _, ti := range teamArr {
		teamMap[ti[0]] = teamCont{ti[1], ti[2]}
	}
	for _, values := range selected {
		eztools.ShowStr(values[len(fieldNames)] + ": ")
		for i := 0; i < len(fieldNames); i++ {
			if i != 2 {
				list1Field(i, values[i])
			}
		}
		eztools.ShowStr("team: " + teamMap[values[2]].str)
		if teamMap[values[2]].leader == values[len(fieldNames)] {
			eztools.ShowStr(", leader")
		} else {
			eztools.ShowStr(", member")
		}
		eztools.ShowStrln("")
	}
	return nil
}
