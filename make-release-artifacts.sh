# Clean dist directory
rm -rf dist
mkdir -p dist

# Get new tags from remote
git fetch --tags

# Get latest tag name
latestTag=$(git describe --tags `git rev-list --tags --max-count=1` | sed -e 's/-rc[0-9][0-9]*//')

# Exclude files without the Apache license header
for i in $(git ls-files); do
   case "$i" in
   (LICENSE|NOTICE|message/common.pb.go|message/requests.pb.go|message/responses.pb.go|test-fixtures/calcite.png|Gopkg.lock|Gopkg.toml);; # add exceptions here
   (*) grep -q "Licensed to the Apache Software Foundation" $i || echo "$i has no header";;
   esac
done

product=apache-calcite-avatica-go
tarFile=$product-src-$latestTag.tar.gz

# Checkout latest tag
git checkout $latestTag

# Make tar
tar -zcvf dist/$tarFile --transform "s/^\./$product-src-$latestTag/g" --exclude "dist" --exclude ".git" .

cd dist

# Calculate MD5
gpg --print-md MD5 $tarFile > $tarFile.md5

# Calculate SHA256
gpg --print-md SHA256 $tarFile > $tarFile.sha256

# Sign
gpg --armor --output $tarFile.asc --detach-sig $tarFile

# End
