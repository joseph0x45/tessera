package models

type App struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type User struct {
	ID       string `db:"id"`
	AppID    string `db:"app_id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
