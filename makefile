# Name
BINARY=wapty
IMPORTPATH=`go list`/

# Variables
VERSION=0.2.0
BUILD=`git rev-parse --short HEAD`

LDFLAGS=-ldflags "-X ${IMPORTPATH}cli.Version=${VERSION} -X ${IMPORTPATH}cli.Commit=${BUILD}"
LDFLAGS_RELEASE=-ldflags "-X ${IMPORTPATH}cli.Version=${VERSION} -X ${IMPORTPATH}cli.Commit=${BUILD} -X ${IMPORTPATH}cli.Build=Release"

.DEFAULT_GOAL: ${BINARY}

# Just build the wapty
# TODO call gopherjs
${BINARY}: buildjs rebind
	# Building the executable.
	go build ${LDFLAGS_RELEASE} -o ${BINARY}

run:
	# This will make rice use data that is on disk, creates a lighter executable
	# and it is faster to build
	-rm ui/rice-box.go >& /dev/null
	# Generating JS
	cd ui/gopherjs/ && gopherjs build -o ../static/gopherjs.js
	# Done generating JS, launching wapty
	go run -race ${LDFLAGS} wapty.go

fast: run

test: buildjs rebind
	go test -race ${LDFLAGS} ./...

testv: buildjs rebind
	go test -v -x -race ${LDFLAGS} ./...

buildjs:
	# Regenerating minified js
	cd ui/gopherjs/ && gopherjs build -m -o ../static/gopherjs.js 
	# Remove mappings
	rm ui/static/gopherjs.js.map

rebind:
	# Cleaning and re-embedding assets
	cd ui && rm rice-box.go 1>/dev/null 2>/dev/null; rice embed-go

install: buildjs rebind
	# Installing the executable
	go install ${LDFLAGS_RELEASE}

installdeps:
	# Installing dependencies to embed assets
	go get github.com/GeertJohan/go.rice/...
	# Installing dependencies to build JS
	go get github.com/gopherjs/gopherjs
	go get github.com/gopherjs/websocket/...
	# Installing Diff dependencies
	go get github.com/fatih/color
	go get github.com/pmezard/go-difflib/difflib

updatedeps:
	# Updating dependencies to embed assets
	go get -u github.com/GeertJohan/go.rice/...
	# Updating dependencies to build JS
	go get -u github.com/gopherjs/gopherjs
	go get -u github.com/gopherjs/websocket/...
	# Updating Diff dependencies
	go get -u github.com/fatih/color
	go get -u github.com/pmezard/go-difflib/difflib

clean:
	# Cleaning all generated files
	-rm ui/rice-box.go
	-rm ui/static/gopherjs.js*
	go clean
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
