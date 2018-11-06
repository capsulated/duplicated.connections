package model

type User struct {
	ID        int64   `reindex:"id,,pk"`
	UserId 	  int64   `reindex:"user_id"`
	Address   string  `reindex:"address"`
	Timestamp string  `reindex:"timestamp"`
}