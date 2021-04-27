package example

//go:generate juryrig gen -o zz.mapper.impl.go

// +juryrig:mapper:MapperImpl
type Mapper interface {
	// +juryrig:link:ef.title->title
	// +juryrig:link:ef.runtime->runtime
	// +juryrig:ignore:director
	// +juryrig:linkfunc:eu->ToInternalUser->user
	ToInternalUserFilm(ef *ExternalFilm, eu *ExternalUser) *InternalUserFilm
	// +juryrig:link:eu.username->username
	ToInternalUser(eu *ExternalUser) *InternalUser
}
