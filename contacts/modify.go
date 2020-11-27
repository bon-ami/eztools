package contacts

import (
	"database/sql"
	"strconv"

	"github.com/bon-ami/eztools"
)

//To modify a record
func Modify(db *sql.DB, id int) error {
	initFieldNames()
	fields := make([]string, 0)
	values := make([]string, 0)
	var field, value []string
	for i := 0; i < len(fieldNames); i++ {
		if i != 2 {
			field, value = prompt1Field(fieldNames[i], nil, nil)
			if len(value) > 0 {
				fields = append(fields, field...)
				values = append(values, value...)
			}
		}
	}
	team, err := eztools.ChoosePairOrAdd(db, eztools.TblTEAM, true)
	if err != nil {
		return err
	}
	fields = append(fields, fieldNames[2])
	values = append(values, strconv.Itoa(team))
	err = eztools.UpdateWtParams(db, eztools.TblCONTACTS,
		eztools.FldID+"="+strconv.Itoa(id), fields, values, false)
	return err
}
