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
   (message/common.pb.go|message/requests.pb.go|message/responses.pb.go|Gopkg.lock|Gopkg.toml);;

   # Binaries
   (test-fixtures/calcite.png);;

   (*) grep -q "Licensed to the Apache Software Foundation" $i || echo "$i has no header";;
   esac
done

tagWithoutRC=$(echo $tag | sed -e 's/-rc[0-9][0-9]*//')
product=apache-calcite-avatica-go
tarFile=$product-src-$tagWithoutRC.tar.gz
releaseDir=$product-$tag

#Make release dir
mkdir -p dist/$releaseDir

# Checkout tag
if ! git checkout $tag; then
    echo "Could not check out tag $tag. Does it exist?"
    exit 1
fi

# Make tar
tar -zcvf dist/$releaseDir/$tarFile --transform "s/^\./$product-src-$tagWithoutRC/g" --exclude "dist" --exclude ".git" .

cd dist/$releaseDir

# Calculate SHA256
gpg --print-md SHA256 $tarFile > $tarFile.sha256

# Sign
gpg --armor --output $tarFile.asc --detach-sig $tarFile

# End
