package main

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/restream/reindexer"
	_ "github.com/restream/reindexer/bindings/builtin"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

// Define struct with reindex tags
type User struct {
	ID        int64   `reindex:"id,,pk"`
	UserId 	  int64   `reindex:"user_id"`
	Address   string  `reindex:"address"`
	Timestamp string  `reindex:"timestamp"`
}

type Response struct {
	Dupes 	bool 	`json:"dupes"`
}

var (
	strContentType = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

var db *reindexer.Reindexer

func main() {
	db = reindexer.NewReindex("builtin:///usersdb")
	db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), User{})

	fmt.Println("Start fill database...")

	// Generate dataset
	for i := 0; i < 50000000; i++ {
		err := db.Upsert("users", &User{
			ID:          int64(i),
			UserId: 	 int64(i),
			Address:     "127.0.0." + strconv.Itoa(i),
			Timestamp: 	time.Now().Format(time.RFC850),
		})
		if err != nil {
			panic(err)
		}
	}

	err := db.Upsert("users", &User{ ID: 1, UserId: 1, Address: "127.0.0.1", Timestamp: 	time.Now().Format(time.RFC850) })
	err = db.Upsert("users", &User{ ID: 2, UserId: 2, Address: "127.0.0.1", Timestamp: 	time.Now().Format(time.RFC850) })
	err = db.Upsert("users", &User{ ID: 3, UserId: 1, Address: "127.0.0.2", Timestamp: 	time.Now().Format(time.RFC850) })
	err = db.Upsert("users", &User{ ID: 4, UserId: 2, Address: "127.0.0.2", Timestamp: 	time.Now().Format(time.RFC850) })
	err = db.Upsert("users", &User{ ID: 5, UserId: 2, Address: "127.0.0.3", Timestamp: 	time.Now().Format(time.RFC850) })
	err = db.Upsert("users", &User{ ID: 6, UserId: 3, Address: "127.0.0.3", Timestamp: 	time.Now().Format(time.RFC850) })
	err = db.Upsert("users", &User{ ID: 7, UserId: 3, Address: "127.0.0.1", Timestamp: 	time.Now().Format(time.RFC850) })
	err = db.Upsert("users", &User{ ID: 8, UserId: 4, Address: "127.0.0.1", Timestamp: 	time.Now().Format(time.RFC850) })

	if err != nil { panic(err) }

	fmt.Println("DB filled")

	router := fasthttprouter.New()
	router.GET("/:user_id_1/:user_id_2", Compare)

	fasthttp.ListenAndServe(":12345", router.Handler)
}

func Compare(ctx *fasthttp.RequestCtx) {
	start := time.Now()
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	u1, err := strconv.Atoi(ctx.UserValue("user_id_1").(string))
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	u2, err := strconv.Atoi(ctx.UserValue("user_id_2").(string))
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	query := db.Query("users").
		WhereInt("user_id", reindexer.SET, u1, u2).
		ReqTotal()

	iterator := query.Exec()
	defer iterator.Close()

	duplicateFrequency := make(map[string]int)

	iter := iterator.Count()
	for iterator.Next() {
		elem := iterator.Object().(*User)
		duplicateFrequency[elem.Address]++
	}

	var res bool
	if iter - len(duplicateFrequency) >= 2 {
		res = true
	} else {
		res = false
	}

	response := &Response{
		Dupes: res,
	}
	body, err := json.Marshal(response)
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	fmt.Printf("%s\n", time.Since(start).String())
}