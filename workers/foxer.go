package workers

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/logiqone/foxed.nesthorn/model"
	"github.com/restream/reindexer"
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

	// f.Db = reindexer.NewReindex("builtin:///usersdb")
	f.Db = reindexer.NewReindex("cproto://127.0.0.1:6534/usersdb")
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
	//for i := 0; i < 1000000; i++ {
	//	err = f.Db.Upsert("users", &model.User{
	//		ID:          int64(i),
	//		UserId: 	 int64(i),
	//		Address:     "127.0.0." + strconv.Itoa(i),
	//		Timestamp: 	time.Now().Format(time.RFC850),
	//	})
	//	if err != nil {
	//		return err
	//	}
	//}

	err = f.Db.Upsert("users", &model.User{ ID: 1, UserId: 1, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 2, UserId: 2, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 3, UserId: 1, Address: "127.0.0.2", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 4, UserId: 2, Address: "127.0.0.2", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 5, UserId: 2, Address: "127.0.0.3", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 6, UserId: 3, Address: "127.0.0.3", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 7, UserId: 3, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 8, UserId: 4, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })

	err = f.Db.Upsert("users", &model.User{ ID: 9, UserId: 3, Address: "127.0.0.7", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 10, UserId: 3, Address: "127.0.0.7", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 11, UserId: 3, Address: "127.0.0.7", Timestamp: time.Now().Format(time.RFC850) })
	err = f.Db.Upsert("users", &model.User{ ID: 12, UserId: 3, Address: "127.0.0.7", Timestamp: time.Now().Format(time.RFC850) })

	if err != nil {
		return err
	}

	fmt.Println("DB filled")
	return nil
}


func (f *Foxer) Handler(ctx *fasthttp.RequestCtx) {
	start := time.Now()
	ctx.Response.Header.SetCanonical(f.StrContentType, f.StrApplicationJSON)

	u1, err := strconv.ParseInt(ctx.UserValue("user_id_1").(string), 10, 64)
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	u2, err := strconv.ParseInt(ctx.UserValue("user_id_2").(string), 10, 64)
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

func (f *Foxer) Compare(u1 int64, u2 int64) bool {
	query := f.Db.Query("users").
		WhereInt64("user_id", reindexer.SET, u1, u2).
		ReqTotal()

	iterator := query.Exec()
	defer iterator.Close()

	u1UniqueIp := make(map[string]struct{})
	u2UniqueIp := make(map[string]struct{})

	for iterator.Next() {
		elem := iterator.Object().(*model.User)
		if elem.UserId == u1 {
			u1UniqueIp[elem.Address] = struct{}{}
		} else {
			u2UniqueIp[elem.Address] = struct{}{}
		}
	}

	frq := 0
	for key := range u1UniqueIp {
		_, ok := u2UniqueIp[key]
		if ok {
			frq++
		}
	}

	fmt.Printf("u1: %v\n", u1UniqueIp)
	fmt.Printf("u2: %v\n", u2UniqueIp)
	fmt.Printf("frq: %v\n", frq)

	var res bool
	if frq >= 2 {
		res = true
	} else {
		res = false
	}

	return res
}