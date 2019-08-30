#/bin/env sh

export GOPATH="${WORKSPACE}/.go"
export GOROOT="${WORKSPACE}/go"
PATH="${GOROOT}/bin:${WORKSPACE}/bin:$PATH"

LATEST=$(curl -s https://golang.org/VERSION?m=text)
mkdir -p "$GOPATH"
mkdir -p "$WORKSPACE/bin"
if [ ! -d "$GOROOT" ] ; then
  # Download the latest stable build
  wget -q https://dl.google.com/go/$LATEST.linux-amd64.tar.gz
  tar xf $LATEST.linux-amd64.tar.gz -C "${WORKSPACE}"
  rm -f $LATEST.linux-amd64.*
fi

# Upgrade if there is one
go get -u golang.org/dl/$LATEST

# UPX (optional)
if [ ! -x "${WORKSPACE}/bin/upx" ] ; then 
  LATEST_UPX=$(curl -s https://github.com/upx/upx/releases/latest | sed -n 's/.*tag\/v\(.*\)\".*/\1/p')
  wget -q https://github.com/upx/upx/releases/download/v${LATEST_UPX}/upx-${LATEST_UPX}-amd64_linux.tar.xz
  tar xf upx-${LATEST_UPX}-amd64_linux.tar.xz --strip=1 -C "${WORKSPACE}/bin" --wildcards '*/upx'
fi
