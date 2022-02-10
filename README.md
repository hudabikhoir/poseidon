![Poseidon](poseidon-logo.png)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=hudabikhoir_poseidon&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=hudabikhoir_poseidon)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=hudabikhoir_poseidon&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=hudabikhoir_poseidon)
[![Open Source? Yes!](https://badgen.net/badge/Open%20Source%20%3F/Yes%21/blue?icon=github)](https://github.com/Naereen/badges/)

# POSEIDON

POSEIDON is a sample REST API build using echo server.

The code implementation was inspired by port and adapter pattern or known as [hexagonal](https://content.octo.com/en/hexagonal-architecture-three-principles-and-an-implementation-example):
In general, Poseidon is divided into 3 major parts, namely primary (driving adapter), business, and secondary (driven adapter).

- **Primary / driving adapter**<br/>driving adapter is a technology that we use to interact with users such as REST API, Graphql, gRPC, and so on. (also called user-side adapters in hexagonal's term)
- **Business**<br/>Contains all the logic in domain business. Also called this as a service. All the interface of repository needed and the implementation of the service itself will be put here.
- **Secondary / driven adapter**<br/>Contains implementations of interfaces defined in the business such as databases, external APIs, clouds, and so on. (also called as server-side adapters in hexagonal's term)

```
.
├── LICENSE
├── README.md
├── api
│   ├── common
│   │   └── dresponse.go
│   ├── insomnia.json
│   ├── router.go
│   └── v1
│       └── content
│           ├── controller.go
│           ├── request
│           │   ├── create_content.go
│           │   └── update_content.go
│           └── response
│               ├── create_new_content.go
│               ├── get_content_by_id.go
│               └── get_content_by_tag.go
├── app
│   ├── main.go
│   └── modules
│       └── modules.go
├── business
│   ├── content
│   │   ├── item.go
│   │   ├── service.go
│   │   ├── service_test.go
│   │   └── spec
│   │       └── upsert_item.go
│   └── error.go
├── config
│   ├── config.go
│   └── config.yaml
├── go.mod
├── go.sum
├── modules
│   └── repository
│       └── content
│           ├── couchdb_repo.go
│           ├── factory.go
│           ├── mongo_repo.go
│           └── mysql_repo.go
├── poseidon.png
├── run.sh
└── util
    ├── dbdriver.go
    └── idgen.go
```

## How To Run Server

### Prerequisite

- Install `reflex` to enable file watcher for hot reload.
  ```console
  $ go get github.com/cespare/reflex
  ```
- Copy config.toml.example to config.toml
  ```console
  $ cp config/config.toml.example config/config.toml
  ```
- adjust your environment configuration in config/config.toml file

### Run Server

Just execute code below in your console

```console
$ ./run.sh
```

## How To Consume The API

There are 4 availables API that ready to use:

- GET `/v1/contents/:id`
- GET `/v1/contents/[tag-name]`
- POST `/v1/contents`
- PUT `/v1/contents`

To make it easier please download [Postman](https://www.postman.com/) app and import [this collection](https://raw.githubusercontent.com/hudabikhoir/poseidon/master/boilerplate-golang-v2.postman_collection.json).

Thank for contributting to this repository. Feel free to contribute this repo, guys :)

<a href = "https://github.com/hudabikhoir/poseidon/graphs/contributors">
<img src = "https://contrib.rocks/image?repo=hudabikhoir/poseidon"/>
</a>
