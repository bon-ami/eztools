package eztools

import (
	"database/sql"

	"github.com/go-ldap/ldap/v3"
)

// Authenticate checks whether authenticated.
// Return value: nil if true.
func Authenticate(db *sql.DB, username, password string) error {
	var (
		conn *ldap.Conn
		err  error
	)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	ip, err := GetPairStr(db, TblCHORE, "AuthIP")
	if err != nil {
		return err
	}
	part1, err := GetPairStr(db, TblCHORE, "AuthRootD")
	if err != nil {
		return err
	}
	part2, err := GetPairStr(db, TblCHORE, "AuthRootN")
	if err != nil {
		return err
	}
	conn, err = ldap.DialURL("ldap://" + ip)
	if err != nil {
		return err
	}
	err = conn.Bind(part1+username+part2, password)
	return err
}
