package contacts

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/bon-ami/eztools"
)

func getIDFromAllExec(db *sql.DB, idOnly bool,
	cri string) (int, [][]string, error) {
	var fields []string
	if idOnly {
		fields = []string{eztools.FldID}
	} else {
		fields = fieldNames
	}
	searched, err := eztools.Search(db,
		eztools.TblCONTACTS, cri, fields, "")
	if err != nil {
		return eztools.InvalidID, nil, err
	}
	if len(searched) == 0 {
		return eztools.InvalidID, nil, eztools.ErrNoValidResults
	}
	ret, err := strconv.Atoi(searched[0][0])
	return ret, searched, err
}

//GetIDFromAllNames gets ID when name contained in nick
func GetIDFromAllNames(db *sql.DB, name string) (int, error) {
	if db == nil || len(name) == 0 {
		return eztools.InvalidID, eztools.ErrInvalidInput
	}
	res, _, err := getIDFromAllExec(db, true,
		eztools.FldNAME+"=\""+name+"\" OR "+
			"number"+"=\""+name+"\" OR "+
			eztools.FldMAIL+"=\""+name+"\" OR "+
			"ldap"+"=\""+name+"\" OR "+
			"uid"+"=\""+name+"\"")
	if err == eztools.ErrNoValidResults {
		res, _, err = getIDFromAllExec(db, true,
			"nick =\""+name+"\" OR "+
				"nick LIKE \""+name+",%\" OR "+
				"nick LIKE \"%,"+name+",%\" OR "+
				"nick LIKE \"%,"+name+"\"")
	}
	return res, err
}

// getFieldFromAllFields gets matched field names when a field value exactly matched
/*func getFieldIndexFromAllFields(fld string, all [][]string) (ret []int) {
	for _, fields := range all {
		for i, field2chk := range fields {
			if field2chk == fld {
				ret = append(ret, i)
			} else {
				ret = append(ret, eztools.InvalidID)
			}
		}
	}
	return
}*/

// getFieldFromAllFields gets matched field names when a field value exactly matched
func getFieldFromAllFields(fld string, all [][]string) []string {
	var ret []string
	for _, fields := range all {
		for i, field2chk := range fields {
			if field2chk == fld {
				ret = append(ret, fieldNames[i])
			} else {
				ret = append(ret, "")
			}
		}
	}
	return ret
}

func chkFieldMail(db *sql.DB, fld string) (int, error) {
	if db == nil || len(fld) == 0 {
		return eztools.InvalidID, eztools.ErrInvalidInput
	}
	const mailSuf = "@fih-foxconn.com"
	if strings.HasSuffix(fld, mailSuf) {
		fld = strings.TrimSuffix(fld, mailSuf)
		id, _, err := getIDFromAllExec(db, true,
			"\""+fld+"\" IN (mail)")
		return id, err
	}
	return eztools.InvalidID, nil
}

func chkFieldOther(db *sql.DB, fld string) (int, [][]string, error) {
	initFieldNames()
	return getIDFromAllExec(db, false,
		"\""+fld+"\" IN ("+fieldNameStr)
}

// GetIDnFieldIndexFromAllFields gets ID and matched field names when a field exactly matched
/*func GetIDnFieldIndexFromAllFields(db *sql.DB, fld string) (int, []int, error) {
	id, err := chkFieldMail(db, fld)
	if err != nil {
		return id, nil, err
	}
	if id != eztools.InvalidID {
		return id, []int{func() int {
			if fieldNames[5] == eztools.FldMAIL {
				return 5
			}
			for i, v := range fieldNames {
				if v == eztools.FldMAIL {
					return i
				}
			}
			return eztools.InvalidID
		}()}, err
	}
	id, searched, err := chkFieldOther(db, fieldNameStr)
	if err == nil && id != eztools.InvalidID {
		f := getFieldIndexFromAllFields(fld, searched)
		if f != nil {
			return id, f, nil
		}
	}
	return id, nil, err
}*/

// GetIDnFieldFromAllFields gets ID and matched field names when a field exactly matched
func GetIDnFieldFromAllFields(db *sql.DB, fld string) (int, []string, error) {
	id, err := chkFieldMail(db, fld)
	if err != nil {
		return id, nil, err
	}
	if id != eztools.InvalidID {
		return id, []string{eztools.FldMAIL}, err
	}
	id, searched, err := chkFieldOther(db, fieldNameStr)
	if err == nil && id != eztools.InvalidID {
		f := getFieldFromAllFields(fld, searched)
		if f != nil {
			return id, f, nil
		}
	}
	return id, nil, err
}

// GetFieldIndexFromAllFieldsByID get matched fields when a field exactly matched by ID
/*func GetFieldIndexFromAllFieldsByID(db *sql.DB, id int, fld string) ([]int, error) {
	initFieldNames()
	searched, err := eztools.Search(db, eztools.TblCONTACTS,
		eztools.FldID+"="+strconv.Itoa(id), fieldNames, "")
	if err != nil {
		return nil, err
	}
	if len(searched) < 1 {
		return nil, eztools.ErrNoValidResults
	}
	return getFieldIndexFromAllFields(fld, searched), nil
}*/

// GetFieldFromAllFieldsByID get matched fields when a field exactly matched by ID
func GetFieldFromAllFieldsByID(db *sql.DB, id int, fld string) ([]string, error) {
	initFieldNames()
	searched, err := eztools.Search(db, eztools.TblCONTACTS,
		eztools.FldID+"="+strconv.Itoa(id), fieldNames, "")
	if err != nil {
		return nil, err
	}
	if len(searched) < 1 {
		return nil, eztools.ErrNoValidResults
	}
	return getFieldFromAllFields(fld, searched), nil
}

func get1Field(db *sql.DB, id, field string) (string, error) {
	searched, err := eztools.Search(db, eztools.TblCONTACTS,
		eztools.FldID+"="+id,
		[]string{field}, "")
	if err != nil {
		return "", err
	}
	return searched[0][0], nil
}

//GetName gets contact name from ID
func GetName(db *sql.DB, id string) (string, error) {
	return get1Field(db, id, eztools.FldNAME)
}

//GetMail gets contact mail from ID
func GetMail(db *sql.DB, id string) (string, error) {
	return get1Field(db, id, eztools.FldMAIL)
}

//GetLdap gets ldap from ID
func GetLdap(db *sql.DB, id string) (string, error) {
	return get1Field(db, id, "ldap")
}
