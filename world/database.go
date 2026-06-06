package world

import (
	_ "modernc.org/sqlite"
	_ "xorm.io/core"
	"xorm.io/xorm"
)

var db *xorm.Engine

func init() {
	var err error
	//db, err = xorm.NewEngine("sqlite3", "file::memory:?cache=shared")
	db, err = xorm.NewEngine("sqlite", NASADataFile)
	if err != nil {
		panic(err)
	}

	// 同步模型到数据库
	err = db.Sync2(new(ObserveData))
	if err != nil {
		panic(err)
	}
	err = db.Sync2(new(CelestialBody))
	if err != nil {
		panic(err)
	}
}
