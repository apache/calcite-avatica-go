@ECHO OFF
rem Licensed to the Apache Software Foundation (ASF) under one or more
rem contributor license agreements.  See the NOTICE file distributed with
rem this work for additional information regarding copyright ownership.
rem The ASF licenses this file to you under the Apache License, Version 2.0
rem (the "License"); you may not use this file except in compliance with
rem the License.  You may obtain a copy of the License at

rem http://www.apache.org/licenses/LICENSE-2.0

rem Unless required by applicable law or agreed to in writing, software
rem distributed under the License is distributed on an "AS IS" BASIS,
rem WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
rem See the License for the specific language governing permissions and
rem limitations under the License.

SET AVATICA_VER=rel/avatica-1.11.0

IF EXIST message\ rmdir /Q /S message
IF EXIST avatica-tmp\ rmdir /Q /S avatica-tmp

git init avatica-tmp
cd avatica-tmp
git remote add origin https://github.com/apache/calcite-avatica/
git config core.sparsecheckout true
echo core/src/main/protobuf/* >> .git/info/sparse-checkout
git fetch --depth=1 origin %AVATICA_VER%
git checkout FETCH_HEAD

cd ..
mkdir message
protoc --proto_path=avatica-tmp/core/src/main/protobuf/ --go_out=import_path=message:message avatica-tmp/core/src/main/protobuf/*.proto

rmdir /Q /S avatica-tmp

echo.
echo Protobufs generated!
