# Easy pgsql go

### Example

```go
package main

import (
	"context"
    "github.com/bendersilver/pgsql"
)

func main() {
	blog.Debug("hi")

	// batch set or update
	var args [][]interface{}
	sql := `INSERT INTO table(id, num) VALUES ($1, $2);`

	for i := 0; i < 100; i++ {
		args = append(args, []interface{}{i, i})
	}
	if err := pgsql.Batch(sql, args); err != nil {
		panic(err)
	}
	args = nil


	// get rows
	var id, val int
	
	sql = `SELECT id, num  FROM table`
	if pool, err := pgsql.PGPool(); err != nil {
		panic(err)
	} else {
		defer pool.Close()
		if conn, err := pool.Acquire(context.Background()); err != nil {
			panic(err)
		} else {
			defer conn.Release()
			rows, _ := conn.Query(context.Background(), sql)
			for rows.Next() {
				h := new(Host)
				err := rows.Scan(&id, &val)
				if err != nil {
					panic(err)
				} else {
					// .. process
				}
			}
			rows.Close()
		}
	}

	// get item
	// pgsql.Get(sql, arg1, arg2 ...) 
	sql = `SELECT id, num  FROM table`
	if err := pgsql.Get(sql).Item(&id, &val); err != nil {
		panic(err)
	}

	// set or update
	// pgsql.Set(sql, arg1, arg2 ...) 
	sql = `UPDATE table SET num = 0`
	if err := pgsql.Set(sql); err != nil {
		panic(err)
	}

}

// initial logging .env file or console
// run myProgramm.bin /path/file/.env
func init() {
	pgsql.Init()
}
```


.env file requery
```
PG_PASS=password
PG_USER=user
PG_DB=db_name
PG_HOST=host
PG_PORT=5372
```