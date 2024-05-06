package schema

type User struct {
	ID        uint64 `db:"id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Password  string `db:"password"`
}

type Book struct {
}
