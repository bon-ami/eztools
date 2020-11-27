package eztools

import (
	"database/sql"
	"errors"
	"strconv"
)

const (
	// TblCHORE contains misc items
	TblCHORE = "chore"
	// TblTOOL contains xTS names
	TblTOOL = "tool"
	// TblVER contains versions of xTS
	TblVER = "ver"
	// TblANDROID contains android versions
	TblANDROID = "android"
	// TblGOOGLE contains xTS requirements
	TblGOOGLE = "google"
	// TblPRODUCT contains product names
	TblPRODUCT = "product"
	// TblPRODGLE contains matches of products and requirements
	TblPRODGLE = "prodgle"
	// TblPRODFO contains matches of products and info
	TblPRODFO = "prodfo"
	// TblBIT contains bit info
	TblBIT = "bit"
	// TblPHASE contains phase info
	TblPHASE = "phase"
	// TblTEAM contains team names
	TblTEAM = "team"
	// TblCONTACTS contains contact info
	TblCONTACTS = "contacts"
	// TblWEEKLYTASKBARS contains bars in report
	TblWEEKLYTASKBARS = "weeklyTaskBars"
	// TblWEEKLYTASKTITLES contains titles in report
	TblWEEKLYTASKTITLES = "weeklyTaskTitles"
	// TblWEEKLYTASKCURR contains titles of current week
	TblWEEKLYTASKCURR = "weeklyTaskCurr"
	// TblWEEKLYTASKNEXT contains titles of next week
	TblWEEKLYTASKNEXT = "weeklyTaskNext"
	// TblWEEKLYTASKDESC contains descriptions of tasks
	TblWEEKLYTASKDESC = "weeklyTaskDesc"
	// TblWEEKLYTASKWORK contains items of tasks
	TblWEEKLYTASKWORK = "weeklyTaskWork"

	// FldID is the name of field id
	FldID = "id"
	// FldSTR is the name of field str
	FldSTR = "str"

	// FldPHASE is the name of field phase
	FldPHASE = "phase"
	// FldTOOL is the name of field tool
	FldTOOL = "tool"
	// FldANDROID is the name of field android
	FldANDROID = "android"
	// FldPRODUCT is the name of field product
	FldPRODUCT = "product"
	// FldGOOGLE is the name of field google
	FldGOOGLE = "google"
	// FldVER is the name of field versions
	FldVER = "ver"
	// FldREQ is the name of field requirements
	FldREQ = "req"
	// FldEXP is the name of field expiry
	FldEXP = "exp"
	// FldBIT is the name of field bit
	FldBIT = "bit"

	// contacts
	// FldLEADER is the name of field leader
	FldLEADER = "leader"
	// FldNUMBER is the name of field number
	FldNUMBER = "number"
	// FldNAME is the name of field name
	FldNAME = "name"
	// FldTEAM is the name of field team
	FldTEAM = "team"
	// FldEXT is the name of field ext
	FldEXT = "ext"
	// FldPHONE is the name of field phone
	FldPHONE = "phone"
	// FldMAIL is the name of field mail
	FldMAIL = "mail"
	// FldLDAP is the name of field ldap
	FldLDAP = "ldap"
	// FldUID is the name of field uid
	FldUID = "uid"
	// FldNICK is the name of field nick
	FldNICK = "nick"

	// FldSECTION is the name of field section
	FldSECTION = "section"
)

func makeWhere(cri string) string {
	if len(cri) > 0 {
		return " WHERE " + cri
	}
	return ""
}

// Rows2Strings returns arrays from rows
// Parameter db not used
func Rows2Strings(db *sql.DB, rows *sql.Rows) (res [][]string, err error) {
	col, err := rows.Columns()
	if err != nil {
		return
	}
	colLen := len(col)
	rawRes := make([][]byte, colLen)
	dest := make([]interface{}, colLen)
	for i := range rawRes {
		dest[i] = &rawRes[i]
	}
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			LogErrPrint(err)
		} else {
			row1 := make([]string, 0)
			for _, raw := range rawRes {
				row1 = append(row1, string(raw))
			}
			res = append(res, row1)
		}
	}
	return
}

func logRows(db *sql.DB, rows *sql.Rows) ([][]string, error) {

	res, err := Rows2Strings(db, rows)
	LogPrint("selected result begins.")
	LogPrint(res)
	LogPrint("selected result ends.")
	return res, err
}

// Search gets values of fields "sel" from "table", using "cri" as "WHERE",
//        with "more" appended to SQL command.
// Parameters: more: will not be prefixed with space automatically
func Search(db *sql.DB, table string, cri string, sel []string, more string) ([][]string, error) {
	var selStr string
	if sel != nil {
		i := 0
		selStr = sel[i]
		for i++; i < len(sel); i++ {
			selStr += "," + sel[i]
		}
	} else {
		selStr = "*"
	}
	if Debugging {
		if Verbose > 1 {
			LogPrint("SELECT " + selStr + " from " + table + makeWhere(cri) + more)
		}
	}
	rows, err := db.Query("SELECT " + selStr + " from " + table + makeWhere(cri) + more)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if Debugging {
		if Verbose > 2 {
			return logRows(db, rows)
		}
	}
	return Rows2Strings(db, rows)
}

// UpdateWtParams updates "fields" in "table" with "values", using "cri" as "WHERE".
// Parameter yes2all = no confirmation in debug mode. Always no confirmation in non-debug mode.
func UpdateWtParams(db *sql.DB, table, cri string, fields, values []string, yes2all bool) error {
	if len(fields) < 1 || len(values) < 1 {
		return ErrInvalidInput
	}
	u := fields[0] + "=\"" + values[0] + "\""
	for i := 1; i < len(fields) && i < len(values); i++ {
		u = u + ", " + fields[i] + "="
		if values[i] != "NULL" {
			u = u + "\""
		}
		u = u + values[i]
		if values[i] != "NULL" {
			u = u + "\""
		}
	}
	if Debugging && !yes2all {
		LogPrint("UPDATE " + table + " SET " + u + makeWhere(cri))
		if !ChkCfmNPrompt("continue", "n") {
			return errors.New("ABORTED")
		}
	}
	_, err := db.Exec("UPDATE " + table + " SET " + u + makeWhere(cri))
	if err != nil {
		LogErrPrint(err)
	}
	return err
}

// AddWtParamsUniq adds "values" to "fields", if no duplicate records in existence.
// Parameter yes2all = no confirmation in debug mode. Always no confirmation in non-debug mode.
func AddWtParamsUniq(db *sql.DB, table string, fields []string, values []string, yes2all bool) (int, error) {
	var cri string
	for i := 0; i < len(fields) && i < len(values); i++ {
		if i != 0 {
			cri += " AND "
		}
		cri += fields[i] + "=" + values[i]
	}
	if len(cri) < 1 {
		return InvalidID, ErrInvalidInput
	}
	//Does ID field exist, to return it when in existence?
	_, err := Search(db, table, "", []string{FldID}, "")
	idExists := true
	var fields2Chk []string
	if err != nil {
		//TODO: check which err it is
		idExists = false
		fields2Chk = make([]string, 1)
		fields2Chk[0] = FldID
	} else {
		fields2Chk = fields
	}
	searched, err := Search(db, table, cri, fields2Chk, "")
	if err != nil {
		if idExists && len(searched) > 0 {
			ret, err := strconv.Atoi(searched[0][0])
			if err == nil {
				return ret, ErrInExistence
			}
		}
		return InvalidID, ErrInExistence
	}
	if len(searched) > 0 {
		return InvalidID, ErrNoValidResults
	}
	return AddWtParams(db, table, fields, values, yes2all)
}

// AddWtParams adds "values" to "fields", no matter whether duplicate records in existence.
// Parameter yes2all = no confirmation in debug mode. Always no confirmation in non-debug mode.
func AddWtParams(db *sql.DB, table string, fields []string, values []string, yes2all bool) (int, error) {
	if len(fields) < 1 || len(values) < 1 {
		return 0, ErrInvalidInput
	}
	f := fields[0]
	v := "\"" + values[0] + "\""
	for i := 1; i < len(fields) && i < len(values); i++ {
		f = f + ", " + fields[i]
		v = v + ", \"" + values[i] + "\""
	}
	if Debugging && !yes2all {
		LogPrint("INSERT INTO " + table + " (" + f + ") VALUES (" + v + ")")
		if !ChkCfmNPrompt("continue", "n") {
			return DefID, errors.New("ABORTED")
		}
	}
	res, err := db.Exec("INSERT INTO " + table + " (" + f + ") VALUES (" + v + ")")
	var id int
	if err != nil {
		LogErrPrint(err)
	} else {
		var num int64
		num, err = res.LastInsertId()
		if err != nil {
			LogErrPrint(err)
		} else {
			id = (int)(num)
		}
	}
	return id, err
}

// DeleteWtParams deletes items with specified WHERE clause
func DeleteWtParams(db *sql.DB, table, where string) error {
	cmd := "DELETE FROM " + table + makeWhere(where)
	if Debugging {
		LogPrint(cmd)
		if !ChkCfmNPrompt("continue", "n") {
			return errors.New("ABORTED")
		}
	}
	_, err := db.Exec(cmd)
	return err
}

// DeleteWtID by ID
func DeleteWtID(db *sql.DB, table, id string) error {
	return DeleteWtParams(db, table, FldID+"="+id)
}
