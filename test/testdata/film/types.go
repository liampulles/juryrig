package film

type ExternalFilm struct {
	title   string
	runtime int
}

type ExternalUser struct {
	username string
	age      int
}

type InternalUser struct {
	username string
}

type InternalUserFilm struct {
	title    string
	runtime  int
	director string
	user     InternalUser
}
