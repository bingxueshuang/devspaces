package db

type User struct {
	Username string
	Password string
	Pubkey   []byte
}

var users []*User

func AddUser(user *User) (ok bool, err error) {
	for _, u := range users {
		if u.Username == user.Username {
			return false, nil
		}
	}
	users = append(users, user)
	return true, nil
}

func MatchUser(user *User) (ok bool, err error) {
	for _, u := range users {
		if u.Username == user.Username {
			return u.Password == user.Password, nil
		}
	}
	return
}

func ListUsers() []User {
	slice := make([]User, 0, len(users))
	for _, u := range users {
		slice = append(slice, *u)
	}
	return slice
}
