GOPATH=/Users/marantz/src/xlax
export GOBIN=${GOPATH}/bin
SPATH=src/github.com/user/xlax

bin/xlax: ${SPATH}/xlax.go
	go install $?
