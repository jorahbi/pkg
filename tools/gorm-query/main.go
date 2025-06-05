package main

import (
	"flag"
	"fmt"

	query "github.com/jorahbi/coco/tools/gorm-query/gen"
)

var dsn = flag.String("dsn", "", "the mysql dsn")
var genPath = flag.String("path", "", "the gen path")
var tables = flag.String("tables", "", "the tables ")
var pName = flag.String("pname", "tools/", "the pkg name ")

type Config struct {
	Db struct {
		Dsn string
	}
}

func main() {
	// t := struct {
	// 	t time.Time
	// }{}
	// fmt.Print(t.t.Local().IsZero())
	flag.Parse()
	if *dsn == "" || *genPath == "" {
		fmt.Println("miss params")
		return
	}
	query.GenQuery(*dsn, *genPath, *tables, *pName)
}

// func build(dsn string) {
// 	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})

// 	q := query.Use(db)
// 	// q.Address
// 	address := q.WithContext(context.Background()).DcomAssetInventoryPlan
// 	condition := []query.DcomAssetInventoryPlanOption{q.DcomAssetInventoryPlan.WithIds([]int32{1, 2})}
// 	mAddress := model.DcomAssetInventoryPlan{}
// 	sql := address.WithOptions(address, condition...).UnderlyingDB().ToSQL(func(tx *gorm.DB) *gorm.DB {
// 		return address.WithOptions(address, condition...).UnderlyingDB().Find(&mAddress)
// 	})
// 	fmt.Println(sql)
// }
