#!/bin/bash

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to you under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Clean dist directory
rm -rf dist
mkdir -p dist

# Get new tags from remote
git fetch --tags

# Prompt for tag to release (defaults to latest tag)
echo -n "Enter tag to release (default: latest tag): "
read tag

if [[ -z $tag ]]; then
    tag=$(git describe --tags `git rev-list --tags --max-count=1`)
    echo "No tag provided. Using the latest tag: $tag"
fi

# Exclude files without the Apache license header
for i in $(git ls-files); do
   case "$i" in
   # The following are excluded from the license header check

   # License files
   (LICENSE|NOTICE);;

   # Generated files
   (message/common.pb.go|message/requests.pb.go|message/responses.pb.go|Gopkg.lock|Gopkg.toml|go.mod|go.sum);;

   # Binaries
   (test-fixtures/calcite.png);;

   (*) grep -q "Licensed to the Apache Software Foundation" $i || echo "$i has no header";;
   esac
done

# Check copyright year in NOTICE
if ! grep -Fq "Copyright 2012-$(date +%Y)" NOTICE; then
    echo "Ending copyright year in NOTICE is not $(date +%Y)"
    exit 1
fi

tagMajorVersion=$(echo $tag | sed -e 's/\..*//')

# Check that go.mod's module path contains the right version
if ! grep -Fq "module github.com/apache/calcite-avatica-go/$tagMajorVersion" go.mod; then
    echo "module declaration in go.mod does not contain the correct version. Expected: $tagMajorVersion"
    exit 1
fi

# Make sure import paths contain the correct version
badImportPaths=false

for i in $(git ls-files); do

    if [[ "$i" == "make-release-artifacts.sh" || "$i" =~ \.md$ ]]; then
        continue
    fi

    lines=$(grep -F '"github.com/apache/calcite-avatica-go' $i)

    if ! [[ -z "$lines" ]]; then
        while read -r line; do
            if ! grep -q "github.com/apache/calcite-avatica-go/$tagMajorVersion" <<< "$line" ; then
                badImportPaths=true
                echo "import for github.com/apache/calcite-avatica-go in $i does not have the correct version ($tagMajorVersion) in its path"
            fi
        done <<< "$lines"
    fi
done

if "$badImportPaths" == true; then
    exit 1
fi

# Check that Avatica versions in both gen-protobuf.bat and gen-protobuf.sh match
EXPECTED_AVATICA_VERSION=$(grep -oP '^export AVATICA_VER="\K[^"]+' gen-protobuf.sh)

if ! grep -Fq "SET AVATICA_VER=$EXPECTED_AVATICA_VERSION" gen-protobuf.bat; then
    echo "AVATICA_VER in gen-protobuf.sh and gen-protobuf.bat does not match. Expected: $EXPECTED_AVATICA_VERSION"
    exit 1
fi

tagWithoutV=$(echo $tag | sed -e 's/v//')
tagWithoutRC=$(echo $tagWithoutV | sed -e 's/-rc[0-9][0-9]*//')
product=apache-calcite-avatica-go
tarFile=$product-$tagWithoutRC-src.tar.gz
releaseDir=dist/$product-$tagWithoutV

#Make release dir
mkdir -p $releaseDir

# Checkout tag
if ! git checkout $tag; then
    echo "Could not check out tag $tag. Does it exist?"
    exit 1
fi

# Make tar
tar -zcf $releaseDir/$tarFile --transform "s/^/$product-$tagWithoutRC-src\//g" $(git ls-files)

cd $releaseDir

# Calculate SHA512
gpg --print-md SHA512 $tarFile > $tarFile.sha512

# Select GPG key for signing
KEYS=()

GPG_COMMAND="gpg2"

get_gpg_keys (){
    GPG_KEYS=$($GPG_COMMAND --list-keys --with-colons --keyid-format LONG)

    KEY_NUM=1

    KEY_DETAILS=""

    while read -r line; do

        IFS=':' read -ra PART <<< "$line"

        if [ ${PART[0]} == "pub" ]; then

            if [ -n "$KEY_DETAILS" ]; then
                KEYS[$KEY_NUM]=$KEY_DETAILS
                KEY_DETAILS=""
                ((KEY_NUM++))

            fi

            KEY_DETAILS=${PART[4]}
        fi

        if [ ${PART[0]} == "uid" ]; then
            KEY_DETAILS="$KEY_DETAILS - ${PART[9]}"
        fi

    done <<< "$GPG_KEYS"

    if [ -n "$KEY_DETAILS" ]; then
        KEYS[$KEY_NUM]=$KEY_DETAILS
    fi
}

get_gpg_keys

if [ "${#KEYS[@]}" -le 0 ]; then
    echo "You do not have any GPG keys available. Exiting..."
    exit 1
fi

echo "You have the following GPG keys:"

for i in "${!KEYS[@]}"; do
        echo "$i) ${KEYS[$i]}"
done

read -p "Select your GPG key for signing: " KEY_INDEX

GPG_KEY=$(sed 's/ -.*//' <<< ${KEYS[$KEY_INDEX]})

if [ -z $GPG_KEY ]; then
    echo "Selected key is invalid. Exiting..."
    exit 1
fi

# Sign
gpg -u $GPG_KEY --armor --output $tarFile.asc --detach-sig $tarFile

echo "Release created!"
# End