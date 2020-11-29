<div align="center"><img src="man-truck.jpg" width="500" alt="Photograph of a man fixing a truck."></div>
<div align="center"><small><sup>A mechanic performing maintenance on a WWII truck.</i></sup></small></div>
<h1 align="center">
  <b>JuryRig</b>
</h1>

<h4 align="center">A tool for generating Go struct mapping code, inspired by mapstruct. </h4>

<p align="center">
  <a href="#status">Status</a> •
  <a href="#run">Run</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

<p align="center">
  <a href="https://github.com/liampulles/juryrig/releases">
    <img src="https://img.shields.io/github/release/liampulles/juryrig.svg" alt="[GitHub release]">
  </a>
  <a href="https://travis-ci.com/liampulles/juryrig">
    <img src="https://travis-ci.com/liampulles/juryrig.svg?branch=master" alt="[Build Status]">
  </a>
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/liampulles/juryrig">
  <a href="https://goreportcard.com/report/github.com/liampulles/juryrig">
    <img src="https://goreportcard.com/badge/github.com/liampulles/juryrig" alt="[Go Report Card]">
  </a>
  <a href="https://codecov.io/gh/liampulles/juryrig">
    <img src="https://codecov.io/gh/liampulles/juryrig/branch/master/graph/badge.svg" />
  </a>
  <a href="https://github.com/liampulles/juryrig/blob/master/LICENSE.md">
    <img src="https://img.shields.io/github/license/liampulles/juryrig.svg" alt="[License]">
  </a>
</p>

## Status

JuryRig is currently in heavy development, and is thus alpha.

## Run

Since JuryRig operates through go generate, once you've added the binary to your `$PATH` you can run

```bash
go generate ./...
```

... or similar.

## Configuration

Given some structs...

```go
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
    title   string
    runtime int
    director string
    user    InternalUser
}
```

We can describe a mapper as follows:

```go
//go:generate juryrig gen -o zz.mapper.impl.go
type Mapper interface {
    // +juryrig:link:ef.title->title
    // +juryrig:link:ef.runtime->runtime
    // +juryrig:ignore:director
    // +juryrig:func:ToInternalUser(eu)->user
    ToInternalUserFilm(ef ExternalFilm, eu ExternalUser) InternalUserFilm
    // +juryrig:link:username->username
    ToInternalUser(eu ExternalUser) InternalUser
}
```

Running go generate will implement the following mapper struct in `zz.mapper.impl.go`:

```go
type MapperImpl struct {}

func NewMapperImpl() *MapperImpl

func (mi *MapperImpl) ToInternalUserFilm(ef ExternalFilm, eu ExternalUser) InternalUserFilm {
    return InternalUserFilm{
        title: ef.title,
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
```

## Contributing

Please submit an issue with your proposal.

## License

See [LICENSE](LICENSE)
