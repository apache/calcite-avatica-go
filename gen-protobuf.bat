SET CALCITE_VER=calcite-1.7.0

rmdir /Q /S message
rmdir /Q /S calcite-tmp

git init calcite-tmp
cd calcite-tmp
git remote add origin https://github.com/apache/calcite/
git config core.sparsecheckout true
echo avatica/core/src/main/protobuf/* >> .git/info/sparse-checkout
git pull --depth=1 origin %CALCITE_VER%

cd ..
mkdir message
protoc --proto_path=calcite-tmp/avatica/core/src/main/protobuf/ --go_out=import_path=message:message calcite-tmp/avatica/core/src/main/protobuf/*.proto

rmdir /Q /S calcite-tmp