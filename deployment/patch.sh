#!/bin/bash

# patch terraform: override the null provider
set -xe


# recompile terraform 'null' provider:
if [ ! -x "${GOPATH}/bin/terraform-provider-null" ] ; then
  CD=$PWD
  # git clone --depth 1 https://github.com/terraform-providers/terraform-provider-null.git "$GOPATH/src/github.com/terraform-providers/terraform-provider-null"
  # cd "$GOPATH/src/github.com/terraform-providers/terraform-provider-null"
  wget -q https://github.com/terraform-providers/terraform-provider-null/archive/master.zip -O null.zip
  mkdir -p "$GOPATH/src/github.com/terraform-providers/"
  unzip -q null.zip -d "$GOPATH/src/github.com/terraform-providers/"
  cd "$GOPATH/src/github.com/terraform-providers/"
  mv terraform-provider-null-master terraform-provider-null
  cd terraform-provider-null
  go install
  cd "$CD"
fi 

cp "$GOPATH/bin/terraform-provider-null" .terraform/plugins/*/terraform-provider-null_*

exit 0
