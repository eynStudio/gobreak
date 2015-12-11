package orm

import (
	"database/sql"
	"log"
)

type TxScope struct {
	*sql.Tx
	Err error
}

func (p *TxScope) Exec(query string, args ...interface{}) (count int64) {
	var r sql.Result
	if r, p.Err = p.Tx.Exec(query, args...); p.Err != nil {
		panic(p.Err)
	}
	if count, p.Err = r.RowsAffected(); p.Err != nil {
		panic(p.Err)
	}
	return
}

func (p *TxScope) Prepare(query string) (stmt *sql.Stmt) {
	if stmt, p.Err = p.Tx.Prepare(query); p.Err != nil {
		panic(p.Err)
	}
	return
}

func (p *TxScope) ExecStmt(stmt *sql.Stmt, args ...interface{}) (count int64) {
	var r sql.Result
	if r, p.Err = stmt.Exec(args...); p.Err != nil {
		panic(p.Err)
	}
	if count, p.Err = r.RowsAffected(); p.Err != nil {
		panic(p.Err)
	}
	return
}

func (p *TxScope) Count(query string, args ...interface{}) (count int64) {
	if p.Err = p.QueryRow(query, args...).Scan(&count); p.Err != nil {
		panic(p.Err)
	}
	return
}

func (p *TxScope) Truncate(table string) *TxScope {
	p.Exec("TRUNCATE TABLE " + table)
	return p
}

func (p *TxScope) Commit() (err error) {
	if e := recover(); e != nil {
		log.Println(e)
		err = p.Tx.Rollback()
	} else {
		err = p.Tx.Commit()
	}
	return
}
