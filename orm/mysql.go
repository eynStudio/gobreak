package orm

import (
	"database/sql"
	"fmt"
	"github.com/eynstudio/gobreak/db/meta"
	"log"
	"strings"
)

type mysql struct {
	commonDialect
}

func (p *mysql) Driver() string { return "mysql" }

func (mysql) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}

func (p *mysql) LoadMeta(db *sql.DB) *meta.MetaDb {
	m := &meta.MetaDb{}
	loadMetaTable(db)
	//	loadMetaTables(db, m)
	log.Println("load mysql meta")
	return m
}

func loadMetaTables(db *sql.DB, m *meta.MetaDb) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		fmt.Println(err)
	}

	var tableName string
	for rows.Next() {
		rows.Scan(&tableName)
		if !strings.Contains(tableName, "(") {
			tbl := meta.MetaTbl{Mc: tableName}
			m.Tbls = append(m.Tbls, tbl)
		}
		//		log.Println(tableName)
	}
	//	for _, it := range m.Tbls {
	//		loadMetaTable(db, &it)
	//	}
}

func loadMetaTable(db *sql.DB) { //}, t *meta.MetaTbl) {
	s := `SHOW COLUMNS IN recruit_sign_up` //+ t.Mc+
	//	s := `select * from information_schema.columns
	//	where table_schema = 'recruit_sign_up'
	//	order by table_name,ordinal_position`
	_, err := db.Query(s) // .Query(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(s)

	//	var mc, lx, nul, pri string
	//	for rows.Next() {
	//		rows.Scan(&mc, &lx, &nul, &pri)
	//		col := meta.MetaCol{Mc: mc, Lx: lx, Nullable: nul == "YES", IsKey: pri == "PRI"}
	//		t.Cols = append(t.Cols, col)
	//		log.Println(t)
	//	}
}
