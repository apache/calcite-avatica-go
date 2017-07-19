SET AVATICA_VER=rel/avatica-1.10.0

rmdir /Q /S message
rmdir /Q /S avatica-tmp

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