package film

type MapperImpl struct{}

func (impl *MapperImpl) ToInternalUserFilm(ef ExternalFilm, eu ExternalUser) InternalUserFilm {
	return InternalUserFilm{
		title:   ef.title,
		runtime: ef.runtime,
		// director: (ignored),
		user: impl.ToInternalUser(eu),
	}
}

func (impl *MapperImpl) ToInternalUser(eu ExternalUser) InternalUser {
	return InternalUser{
		username: eu.username,
	}
}
