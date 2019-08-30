#/bin/env sh

export GOPATH="${WORKSPACE}"
export GOROOT="${GOPATH}/go"
PATH="${GOROOT}/bin:$PATH"

cd "${WORKSPACE}"
LATEST=$(curl https://golang.org/VERSION?m=text)
if ! -d "$GOROOT" ; then
  # Download the latest stable build
  wget https://dl.google.com/go/$LATEST.linux-amd64.tar.gz
  tar xf $LATEST.linux-amd64.tar.gz
fi

# Upgrade if there is one
go get -u golang.org/dl/$LATEST
