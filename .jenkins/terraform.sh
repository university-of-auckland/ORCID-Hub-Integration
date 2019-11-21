#/bin/env sh
PATH="${GOROOT}/bin:${WORKSPACE}/bin:$PATH"

which terraform >/dev/null
if [ $? -ne 0 ] || (terraform  version  | grep -q 'You can update') ; then
  # wget -O terraform.zip -q $(curl -s https://checkpoint-api.hashicorp.com/v1/check/terraform | sed 's/.*current_version":"\([^"]*\)".*current_download_url":"\([^"]*\).*/\2terraform_\1_linux_amd64.zip/')
  VERSION="0.12.11"
  wget -O terraform.zip -q https://releases.hashicorp.com/terraform/${VERSION}/terraform_${VERSION}_linux_amd64.zip
  unzip -qq -o terraform.zip -d "$WORKSPACE/bin"
  rm -f terraform.zip
fi


