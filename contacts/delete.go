package contacts

import (
	"database/sql"
	"strconv"

	"github.com/bon-ami/eztools"
)

func Delete(db *sql.DB, id int) error {
	err := eztools.DeleteWtID(db, eztools.TblCONTACTS,
		strconv.Itoa(id))
	return err
}
