package workers

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/logiqone/foxed.nesthorn/model"
	"github.com/restream/reindexer"
	_ "github.com/restream/reindexer/bindings/builtin"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

type Foxer struct{
	StrContentType 		[]byte
	StrApplicationJSON  []byte
	Db 					*reindexer.Reindexer
}

func (f *Foxer) Init() error {
	f.StrContentType = []byte("Content-Type")
	f.StrApplicationJSON = []byte("application/json")

	f.Db = reindexer.NewReindex("builtin:///usersdb")
	return f.Db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.User{})
}

func (f *Foxer) Start()  error {
	router := fasthttprouter.New()
	router.GET("/:user_id_1/:user_id_2", f.Handler)

	return fasthttp.ListenAndServe(":12345", router.Handler)
}

func (f *Foxer) InitDataFill() (err error) {
	fmt.Println("Start fill database...")

	// Generate dataset
	for i := 0; i < 5000000; i++ {
		err = f.Db.Upsert("users", &model.User{
			ID:          int64(i),
			UserId: 	 int64(i),
			Address:     "127.0.0." + strconv.Itoa(i),
			Timestamp: 	time.Now().Format(time.RFC850),
		})
		if err != nil {
			return err
		}
	}

	err = f.Db.Upsert("users", &model.User{ ID: 1, UserId: 1, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 2, UserId: 2, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 3, UserId: 1, Address: "127.0.0.2", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 4, UserId: 2, Address: "127.0.0.2", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 5, UserId: 2, Address: "127.0.0.3", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 6, UserId: 3, Address: "127.0.0.3", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 7, UserId: 3, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 8, UserId: 4, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })

	if err != nil {
		return err
	}

	fmt.Println("DB filled")
	return nil
}


func (f *Foxer) Handler(ctx *fasthttp.RequestCtx) {
	start := time.Now()
	ctx.Response.Header.SetCanonical(f.StrContentType, f.StrApplicationJSON)

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

	res := f.Compare(u1, u2)

	response := &model.Response{
		Dupes: res,
	}

	body, err := response.MarshalJSON()
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	fmt.Printf("%s\n", time.Since(start).String())
}

func (f *Foxer) Compare(u1 int, u2 int) bool {
	query := f.Db.Query("users").
		WhereInt("user_id", reindexer.SET, u1, u2).
		ReqTotal()

	iterator := query.Exec()
	defer iterator.Close()

	duplicateFrequency := make(map[string]int)

	iter := iterator.Count()
	for iterator.Next() {
		elem := iterator.Object().(*model.User)
		duplicateFrequency[elem.Address]++
	}

	var res bool
	if iter - len(duplicateFrequency) >= 2 {
		res = true
	} else {
		res = false
	}

	return res
}