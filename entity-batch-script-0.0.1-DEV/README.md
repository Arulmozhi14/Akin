# entity-batch-script

entity-batch-script is a script to be executed by the entity. It copies the "users" table of the entity to the entity-api database.

### Installation

This was developed on `go version go1.12.6 darwin/amd64`

Clone it inside your go workspace (example: /home/ubuntu/go/src/)

```sh
$ git clone https://gitlab.ubx.ph/identity/entity-batch-script.git
```

Install the dependencies.

```sh
$ cd entity-batch-script
$ go get
```

Copy the sample config and update as necessary

```sh
$ cp sample.env .env
$ vim .env
```

Running the script

```sh
$ cd entity-batch-script
$ go run main.go
```

Build script into a binary file

```sh
$ cd entity-batch-script
$ go install
```

Execute `bin/entity-batch-script`
```sh
$ bin/entity-batch-script
```
