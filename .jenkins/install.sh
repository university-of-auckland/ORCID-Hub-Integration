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

go get \
  golang.org/x/tools/cmd/cover \
  github.com/mattn/goveralls \
  golang.org/x/lint/golint \
  gotest.tools/gotestsum \
  honnef.co/go/tools/cmd/staticcheck

# Upgrade if there is one
[ "$(go version | cut -d' ' -f3)" != "$LATEST" ] && go get -u golang.org/dl/$LATEST

# UPX (optional)
if [ ! -x "${WORKSPACE}/bin/upx" ] ; then 
  LATEST_UPX=$(curl -s https://github.com/upx/upx/releases/latest | sed -n 's/.*tag\/v\(.*\)\".*/\1/p')
  wget -q https://github.com/upx/upx/releases/download/v${LATEST_UPX}/upx-${LATEST_UPX}-amd64_linux.tar.xz
  tar xf upx-${LATEST_UPX}-amd64_linux.tar.xz --strip=1 -C "${WORKSPACE}/bin" --wildcards '*/upx'
fi

# zip
# if [ ! -x "${WORKSPACE}/bin/zipit" ] ; then 
#   go build -o "${WORKSPACE}/bin/zipit" .jenkins/zipit.go
# fi

# recompile terraform 'null' provider:
if [ ! -x "${GOPATH}/bin/terraform-provider-null" ] ; then
  git clone --depth 1 git@github.com:terraform-providers/terraform-provider-null "$GOPATH/src/github.com/terraform-providers/terraform-provider-null"
  cd "$GOPATH/src/github.com/terraform-providers/terraform-provider-null"

  # wget https://github.com/terraform-providers/terraform-provider-null/archive/master.zip -O null.zip
  # unzip -q null.zip -d "$GOPATH/src/github.com/terraform-providers/"
  # cd "$GOPATH/src/github.com/terraform-providers/"
  # mv terraform-provider-null-master terraform-provider-null
  # cd terraform-provider-null
  go install

fi 
