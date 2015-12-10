package orm

import (
	"database/sql"
)

type TxScope struct {
	*sql.Tx
	Err error
}

func (p *TxScope) Exec(query string, args ...interface{}) (count int64) {
	var r sql.Result
	if p.Err == nil {
		if r, p.Err = p.Tx.Exec(query, args); p.Err == nil {
			count, p.Err = r.RowsAffected()
		}
	}
	return
}

func (p *TxScope) Prepare(query string) (stmt *sql.Stmt) {
	stmt, p.Err = p.Tx.Prepare(query)
	return
}

func (p *TxScope) ExecStmt(stmt *sql.Stmt, args ...interface{}) (count int64) {
	var r sql.Result
	if p.Err == nil {
		if r, p.Err = stmt.Exec(args); p.Err == nil {
			count, p.Err = r.RowsAffected()
		}
	}
	return
}

func (p *TxScope) Count(query string, args ...interface{}) (count int64) {
	p.Err = p.QueryRow(query, args...).Scan(&count)
	return
}
