package user

type User struct {
	Id       uint64
	Name     string
	Email    string
	Password []byte
}

type UserWithoutPassword struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
