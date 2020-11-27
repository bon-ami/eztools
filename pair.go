package eztools

import (
	"database/sql"
	"sort"
	"strconv"
)

// pairs is an int-string pair
type pairs struct {
	id  int
	str string
}

// PairsStr is a collection of pairs
type PairsStr struct {
	pair []pairs
	curr int
}

// pairsI is an int-int pair
type pairsI struct {
	id  int
	val int
}

// PairsInt is a collection of pairsI
type PairsInt struct {
	pair []pairsI
	curr int
}

// GetSortedPairsStr returns all sorted PairsStr from input table
func GetSortedPairsStr(db *sql.DB, table string) (ps *PairsStr, err error) {
	ps = new(PairsStr)
	var (
		id  int
		val string
	)
	err = selectSortedPairs(db, table, &id, &val,
		func(err error) {
			if err != nil {
				LogErr(err)
				return
			}
			ps.pair = append(ps.pair, pairs{id, val})
		})
	return
}

// Next moves to the next item in PairsInt and return the values
// error=ErrOutOfBound when the end met
func (ps *PairsStr) Next() (int, string, error) {
	if ps.curr < len(ps.pair) {
		ps.curr++
		return ps.pair[ps.curr-1].id, ps.pair[ps.curr-1].str, nil
	}
	return 0, "", ErrOutOfBound
}

// Rewind sets the cursor to the beginning of PairsInt
func (ps *PairsStr) Rewind() {
	ps.curr = 0
}

// FindStr find the first ID where str matches input
// return value: ErrNoValidResults when none found
func (ps *PairsStr) FindStr(s string) (int, error) {
	for _, v := range ps.pair {
		if v.str == s {
			return v.id, nil
		}
	}
	return InvalidID, ErrNoValidResults
}

type selectSortedPairsFunc func(error)

func selectSortedPairs(db *sql.DB, table string, id, val interface{}, fn selectSortedPairsFunc) (err error) {
	cmd := "SELECT " + FldID + "," + FldSTR + " from " + table + " order by " + FldID
	if Debugging {
		if Verbose > 2 {
			LogPrint(cmd)
		}
	}
	rows, err := db.Query(cmd)
	if err != nil {
		return
	}
	defer rows.Close()
	err = ErrNoValidResults
	for rows.Next() {
		err = rows.Scan(id, val)
		fn(err)
	}
	return
}

// GetSortedPairsInt returns all sorted PairsInt from input table
func GetSortedPairsInt(db *sql.DB, table string) (pi *PairsInt, err error) {
	pi = new(PairsInt)
	var id, val int
	err = selectSortedPairs(db, table, &id, &val,
		func(err error) {
			if err != nil {
				LogErr(err)
				return
			}
			pi.pair = append(pi.pair, pairsI{id, val})
		})
	return
}

// Next moves to the next item in PairsInt and return the values
// error=ErrOutOfBound when the end met
func (pi *PairsInt) Next() (int, int, error) {
	if pi.curr < len(pi.pair) {
		pi.curr++
		return pi.pair[pi.curr-1].id, pi.pair[pi.curr-1].val, nil
	}
	return 0, 0, ErrOutOfBound
}

// Rewind sets the cursor to the beginning of PairsInt
func (pi *PairsInt) Rewind() {
	pi.curr = 0
}

// GetPairIDFromInt gets ID (int) from value (int)
// When multiple results got, the first one will be taken.
func GetPairIDFromInt(db *sql.DB, table string, val int) (int, error) {
	return GetPairID(db, table, strconv.Itoa(val))
}

// GetPairID gets ID (int) from value (string)
// When multiple results got, the first one will be taken.
func GetPairID(db *sql.DB, table string, str string) (int, error) {
	ret, err := GetPair(db, table, str, FldSTR, FldID)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(ret)
}

// GetPairInt gets value (int) from ID (string)
func GetPairInt(db *sql.DB, table string, id string) (int, error) {
	res, err := GetPairStr(db, table, id)
	if err != nil {
		return InvalidID, err
	}
	return strconv.Atoi(res)
}

// GetPairStrFromInt gets value (string) from ID (int)
func GetPairStrFromInt(db *sql.DB, table string, id int) (string, error) {
	return GetPairStr(db, table, strconv.Itoa(id))
}

// GetPairStr gets value (string) from ID (string)
func GetPairStr(db *sql.DB, table string, id string) (string, error) {
	return GetPair(db, table, id, FldID, FldSTR)
}

// GetPair gets "to" field from "from" field in "table"
// When multiple results got, the first one will be taken.
// return value error = from db.Query;
//                      ErrNoValidResults when no valid results got,
//                                        and LogErrPrint will be called.
func GetPair(db *sql.DB, table, id, from, to string) (string, error) {
	//if Debugging {
	//LogPrint("SELECT " + FldSTR + " from " + table + " where " + FldID + "=" + id)
	//}
	cmd := "SELECT " + to + " from " + table + " where " + from + "=\"" + id + "\""
	if Debugging {
		if Verbose > 2 {
			LogPrint(cmd)
		}
	}
	rows, err := db.Query(cmd)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var str string
	if Debugging && Verbose > 2 {
		var res [][]string
		res, err = logRows(db, rows)
		if err == nil {
			if res != nil {
				str = res[0][0]
			} else {
				err = ErrNoValidResults
			}
		}
	} else {
		err = ErrNoValidResults
		for rows.Next() {
			err = rows.Scan(&str)
			break
		}
	}
	//if err != nil {
	//LogErrPrint(err)
	//}
	return str, err
}

// Locate gets ID (int) from value (string)
// Deprecated: Use GetPairID instead.
func Locate(db *sql.DB, table string, str string) (id int, err error) {
	var (
		rows *sql.Rows
	)
	rows, err = db.Query("SELECT " + FldID + " from " + table + " where " + FldSTR + "=\"" + str + "\" order by " + FldID)
	if err != nil {
		return DefID, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			LogErrPrint(err)
		}
	}
	return id, nil
}

// selectPairs get all values in a table.
func selectPairs(db *sql.DB, table string) (arr []pairs, err error) {
	var (
		id  int
		str string
		l   int
	)
	//put ordered results into [id1, "str1,str2,..."]
	err = selectSortedPairs(db, table, &id, &str,
		func(err error) {
			if err != nil {
				LogErr(err)
				return
			}
			l = len(arr)
			if l > 0 && arr[l-1].id == id {
				arr[l-1].str += ",str"
			} else {
				arr = append(arr, pairs{id, str})
			}
		})
	return
}

// choosePairOrAddExec asks user to choose from a table by ID - value pairs
//                     and allows user to add one, returning the ID user has chosen
func choosePairOrAddExec(db *sql.DB, table string, add bool, uniq bool) (int, error) {
	tblAllAllowed := make(map[string]bool)
	tblAllAllowed[TblANDROID] = true
	arr, err := selectPairs(db, table)
	if err != nil {
		return InvalidID, err
	}
	if add {
		if tblAllAllowed[table] {
			arr = append(arr, pairs{AllID, " <Affect All Items>"})
		}
		arr = append(arr, pairs{DefID, "(default) <Add An Item>"})
	}
	res := choosePairs(arr, "Choose for "+table)
	if !add || res != DefID || res == InvalidID {
		return res, nil
	}
	return addPairWtPrompt(db, table, uniq)
}

// ChoosePair asks user to choose from a table by ID - value pairs.
func ChoosePair(db *sql.DB, table string) (int, error) {
	return choosePairOrAddExec(db, table, false, false)
}

// ChoosePairOrAdd asks user to choose from a table by ID - value pairs,
//            allowing to add one new.
func ChoosePairOrAdd(db *sql.DB, table string, uniq bool) (int, error) {
	return choosePairOrAddExec(db, table, true, uniq)
}

// ChoosePairNAddAssociated asks user to choose by idTable - strTable pairs,
//                          where idTable.str == strTable.id (return value)
// Parameter uniq is of no use
func ChoosePairNAddAssociated(db *sql.DB, idTable, strTable string, uniq bool) (int, error) {
	arr, err := selectPairs(db, idTable)
	if err != nil {
		return InvalidID, err
	}
	arrStr, err := selectPairs(db, strTable)
	if err != nil {
		return InvalidID, err
	}
	sort.Slice(arrStr, func(i, j int) bool {
		return arrStr[i].id < arrStr[j].id
	})
	lenStr := len(arrStr)
	for i := range arr {
		ind, err := strconv.Atoi(arr[i].str)
		if err != nil {
			LogErrPrint(err)
			continue
		}
		i := sort.Search(lenStr, func(i int) bool {
			return ind <= arrStr[i].id
		})
		if i < len(arrStr) && arrStr[i].id == ind {
			arr[i].str = arrStr[i].str
		}
	}
	//if add {
	arr = append(arr, pairs{DefID, "(default) <Add An Item>"})
	//}
	res := choosePairs(arr, "Choose for "+idTable+"+"+strTable)
	if /*!add || */ res != DefID || res == InvalidID {
		return res, nil
	}
	return res, nil //addPairWtPrompt(db, table, uniq)
}

// addPairWtPrompt asks user to create an ID - value pair and returns the ID
// Parameter uniq = ID will be auto generated.
func addPairWtPrompt(db *sql.DB, table string, uniq bool) (int, error) {
	var (
		value, id string
		idInt     int
		err       error
	)
	if !uniq {
		id = PromptStr("Enter ID for " + table + "(" + strconv.Itoa(DefID) + " for new item; ID in existence to add alias): ")
		if len(id) < 1 {
			return DefID, ErrInvalidInput
		}
		idInt, err = strconv.Atoi(id)
		if err != nil {
			ShowStrln("Invalid ID. Default taken.")
			idInt = DefID
		} else {
			if idInt == AllID {
				PromptStr(id + " is used for all items. Probably wrong input!")
			}
		}
	}
	value = PromptStr("Enter value for " + table + "(no default value): ")
	if len(value) < 1 {
		return DefID, ErrInvalidInput
	}
	fields := []string{FldSTR}
	values := []string{value}
	if !uniq {
		if idInt == DefID {
			idInt, err = getMaxID(db, table)
			if err != nil || idInt == DefID {
				//no max ID found
				idInt = defValidID
			} else {
				idInt++
			}
			id = strconv.Itoa(idInt)
		}
		fields = append(fields, FldID)
		values = append(values, id)
	}
	val, err := AddWtParams(db, table, fields, values, false)
	if err != nil {
		return DefID, err
	}
	if !uniq {
		return idInt, err
	}
	return val, err
}

// getMaxID gets the largest ID in a table.
func getMaxID(db *sql.DB, table string) (int, error) {
	rows, err := db.Query("SELECT MAX(" + FldID + ") from " + table)
	if err != nil {
		//LogErrPrint(err)
		return DefID, err
	}
	defer rows.Close()
	rows.Next()
	var id int
	if err = rows.Scan(&id); err != nil {
		return InvalidID, err
	}
	return id, nil
}

// AddPairNoID adds an item with value, where ID should be auto generated
func AddPairNoID(db *sql.DB, table string, str string) (int, error) {
	return AddWtParams(db, table, []string{FldSTR}, []string{str}, false)
}

// AddPair adds an item with ID and value
func AddPair(db *sql.DB, table string, id int, str string) (res int, err error) {
	cmd := "INSERT INTO " + table + " (" + FldID + ", " + FldSTR + ") VALUES (\"" + strconv.Itoa(id) + "\", \"" + str + "\")"
	if Debugging {
		if Verbose > 2 {
			LogPrint(cmd)
		}
	}
	_, err = db.Exec(cmd)
	//if err != nil {
	//LogErrPrint(err)
	//}
	res = id
	return
}

// UpdatePairWtParams updates value by ID
func UpdatePairWtParams(db *sql.DB, table string, id, str string) error {
	cmd := "UPDATE " + table + " SET " + FldSTR + "=\"" +
		str + "\" WHERE " + FldID + "=\"" + id + "\""
	if Debugging {
		if Verbose > 2 {
			LogPrint(cmd)
		}
	}
	_, err := db.Exec(cmd)
	//if err != nil {
	//LogErrPrint(err)
	//}
	return err
}

// UpdatePairID updates ID
func UpdatePairID(db *sql.DB, table string, idOld, idNew string) error {
	cmd := "UPDATE " + table + " SET " + FldID + "=\"" +
		idNew + "\" WHERE " + FldID + "=\"" + idOld + "\""
	if Debugging {
		if Verbose > 2 {
			LogPrint(cmd)
		}
	}
	_, err := db.Exec(cmd)
	//if err != nil {
	//LogErrPrint(err)
	//}
	return err
}
