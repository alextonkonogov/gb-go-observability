package repository

type motivation struct {
	Id      int    `db:"id"`
	Content string `db:"content"`
	Author  string `db:"author"`
	UserId  int    `db:"user_id"`
}
