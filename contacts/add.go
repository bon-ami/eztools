package contacts

import (
	"database/sql"
	"strconv"

	"github.com/bon-ami/eztools"
)

func prompt1Field(field string, fields, values []string) ([]string, []string) {
	const forContact = " for this contact:"
	val := eztools.PromptStr(field + forContact)
	if len(val) > 0 {
		return append(fields, field), append(values, val)
	}
	return fields, values
}

func Add(db *sql.DB, teamI int) (id, teamO int, err error) {
	initFieldNames()
	fields := make([]string, 0)
	values := make([]string, 0)
	for i := 0; i < len(fieldNames); i++ {
		if i != 2 {
			fields, values = prompt1Field(fieldNames[i], fields, values)
		}
	}
	if teamI == eztools.InvalidID {
		teamI, err = eztools.ChoosePairOrAdd(db, eztools.TblTEAM, true)
		if err != nil {
			return -1, -1, err
		}
	}
	fields = append(fields, fieldNames[2])
	values = append(values, strconv.Itoa(teamI))
	id, err = eztools.AddWtParams(db, eztools.TblCONTACTS, fields, values, true)
	return id, teamI, err
}
