package schema

type User struct {
	ID        uint64 `db:"id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Password  string `db:"password"`
}

type Book struct {
	ID              uint64  `db:"id"`
	BelongsToUserID uint64  `db:"belongs_to_user_id"`
	Name            string  `db:"name"`
	Author          string  `db:"author"`
	Genre           string  `db:"genre"`
	Description     string  `db:"description"`
	Latitude        float64 `db:"latitude"`
	Longitude       float64 `db:"longitude"`
}
