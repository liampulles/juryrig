package film

type MapperImpl struct{}

func NewMapperImpl() *MapperImpl

func (mi *MapperImpl) ToInternalUserFilm(ef ExternalFilm, eu ExternalUser) InternalUserFilm {
	return InternalUserFilm{
		title:   ef.title,
		runtime: ef.runtime,
		// director: (ignored)
		user: mi.ToInternalUser(eu),
	}
}

func (mi *MapperImpl) ToInternalUser(eu ExternalUser) InternalUser {
	return InternalUser{
		username: eu.username,
	}
}
