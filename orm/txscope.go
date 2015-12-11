package orm

import (
	"database/sql"
	"log"
)

type TxScope struct {
	*sql.Tx
}

func (p *TxScope) Exec(query string, args ...interface{}) int64 {
	if r, err := p.Tx.Exec(query, args...); err != nil {
		panic(err)
	} else {
		return p.getAffectedRows(r)
	}
}

func (p *TxScope) Prepare(query string) *sql.Stmt {
	if stmt, err := p.Tx.Prepare(query); err != nil {
		panic(err)
	} else {
		return stmt
	}
}

func (p *TxScope) ExecStmt(stmt *sql.Stmt, args ...interface{}) int64 {
	if r, err := stmt.Exec(args...); err != nil {
		panic(err)
	} else {
		return p.getAffectedRows(r)
	}
}

func (p *TxScope) getAffectedRows(r sql.Result) int64 {
	if count, err := r.RowsAffected(); err != nil {
		panic(err)
	} else {
		return count
	}
}
func (p *TxScope) Count(query string, args ...interface{}) (count int64) {
	if err := p.QueryRow(query, args...).Scan(&count); err != nil {
		panic(err)
	}
	return
}

func (p *TxScope) Truncate(table string) int64 {
	return p.Exec("TRUNCATE TABLE " + table)
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
