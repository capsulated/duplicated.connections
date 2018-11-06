package workers_test

import (
	"github.com/logiqone/foxed.nesthorn/model"
	"github.com/logiqone/foxed.nesthorn/workers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Foxer", func() {

	fox := new(workers.Foxer)
	if err := fox.Init(); err != nil {
		Fail("Can't initialise worker!")
	}

	var err error
	err = fox.Db.Upsert("users", &model.User{ ID: 1, UserId: 1, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	if err != nil {
		Fail("Can't Upsert!")
	}
	err = fox.Db.Upsert("users", &model.User{ ID: 2, UserId: 2, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	if err != nil {
		Fail("Can't Upsert!")
	}
	err = fox.Db.Upsert("users", &model.User{ ID: 3, UserId: 1, Address: "127.0.0.2", Timestamp: time.Now().Format(time.RFC850) })
	if err != nil {
		Fail("Can't Upsert!")
	}
	err = fox.Db.Upsert("users", &model.User{ ID: 4, UserId: 2, Address: "127.0.0.2", Timestamp: time.Now().Format(time.RFC850) })
	if err != nil {
		Fail("Can't Upsert!")
	}
	err = fox.Db.Upsert("users", &model.User{ ID: 5, UserId: 2, Address: "127.0.0.3", Timestamp: time.Now().Format(time.RFC850) })
	if err != nil {
		Fail("Can't Upsert!")
	}
	err = fox.Db.Upsert("users", &model.User{ ID: 6, UserId: 3, Address: "127.0.0.3", Timestamp: time.Now().Format(time.RFC850) })
	if err != nil {
		Fail("Can't Upsert!")
	}
	err = fox.Db.Upsert("users", &model.User{ ID: 7, UserId: 3, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })
	if err != nil {
		Fail("Can't Upsert!")
	}
	err = fox.Db.Upsert("users", &model.User{ ID: 8, UserId: 4, Address: "127.0.0.1", Timestamp: time.Now().Format(time.RFC850) })

	Describe("Compare user_id's", func() {
		// Выполняем GET запрос: http://localhost:12345/1/2
		It("It should be duplicates: true", func() {
			Ω(fox.Compare(1, 2)).To(Equal(true))
		})

		// Выполняем GET запрос: http://localhost:12345/1/3
		It("It should be duplicates: false", func() {
			Ω(fox.Compare(1, 3)).To(Equal(false))
		})

		// Выполняем GET запрос: http://localhost:12345/2/1
		It("It should be duplicates: true", func() {
			Ω(fox.Compare(2, 1)).To(Equal(true))
		})

		// Выполняем GET запрос: http://localhost:12345/2/3
		It("It should be duplicates: true", func() {
			Ω(fox.Compare(2, 3)).To(Equal(true))
		})

		// Выполняем GET запрос: http://localhost:12345/3/2
		It("It should be duplicates: false", func() {
			Ω(fox.Compare(3, 2)).To(Equal(true))
		})

		// Выполняем GET запрос: http://localhost:12345/1/4
		It("It should be duplicates: false", func() {
			Ω(fox.Compare(1, 4)).To(Equal(false))
		})

		// Выполняем GET запрос: http://localhost:12345/3/1
		It("It should be duplicates: false", func() {
			Ω(fox.Compare(3, 1)).To(Equal(false))
		})
	})
})
