package usernamedb

type UsernameDB struct {
	userMap       map[string]string
	Loader        DBLoader
	AddUnresolved bool
}

func (u *UsernameDB) Load() (err error) {
	u.userMap, err = u.Loader.Load()
	return
}

func (u *UsernameDB) ToUsername(email string) (string, bool) {
	username, ok := u.userMap[email]
	return "@" + username, ok
}

func (u *UsernameDB) ToUsernames(emails []string) (usernames []string) {
	for _, email := range emails {
		if username, ok := u.ToUsername(email); ok {
			usernames = append(usernames, username)
		} else if u.AddUnresolved {
			usernames = append(usernames, email)
		}
	}
	return
}
