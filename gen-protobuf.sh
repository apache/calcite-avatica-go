#!/usr/bin/env bash

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
set -e

rm -rf message avatica-tmp

export AVATICA_VER="rel/avatica-1.11.0"

mkdir -p avatica-tmp
pushd avatica-tmp &> /dev/null

git init
git remote add origin https://github.com/apache/calcite-avatica.git
git config core.sparsecheckout true
echo "core/src/main/protobuf/*" >> .git/info/sparse-checkout

git fetch --depth=1 origin $AVATICA_VER
git checkout FETCH_HEAD

popd &> /dev/null

mkdir -p message

protoc --proto_path=avatica-tmp/core/src/main/protobuf --go_out=import_path=message:message avatica-tmp/core/src/main/protobuf/*.proto

rm -rf avatica-tmp

echo -e "\nProtobufs generated!"
