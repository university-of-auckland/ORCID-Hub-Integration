#/bin/env sh

export GOPATH="${WORKSPACE}/.go"
export GOROOT="${WORKSPACE}/go"
PATH="${GOROOT}/bin:$PATH"

LATEST=$(curl https://golang.org/VERSION?m=text)
mkdir -p "$GOPATH"
if ! -d "$GOROOT" ; then
  # Download the latest stable build
  wget https://dl.google.com/go/$LATEST.linux-amd64.tar.gz
  tar xf $LATEST.linux-amd64.tar.gz
fi

# Upgrade if there is one
go get -u golang.org/dl/$LATEST
