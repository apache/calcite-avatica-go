# Clean dist directory
rm -rf dist
mkdir -p dist

# Get new tags from remote
git fetch --tags

# Get latest tag name
latestTag=$(git describe --tags `git rev-list --tags --max-count=1`)

# Checkout latest tag
git checkout $latestTag

# Make tar
tar -zcvf dist/calcite-avatica-go-src-$latestTag.tar.gz --transform "s/^\./calcite-avatica-go-src-$latestTag/g" --exclude "dist" .

cd dist

# Calculate MD5
gpg --print-md MD5 calcite-avatica-go-src-$latestTag.tar.gz > calcite-avatica-go-src-$latestTag.tar.gz.md5

# Calculate SHA256
gpg --print-md SHA256 calcite-avatica-go-src-$latestTag.tar.gz > calcite-avatica-go-src-$latestTag.tar.gz.sha256

# Sign
gpg --armor --output calcite-avatica-go-src-$latestTag.gz.asc --detach-sig calcite-avatica-go-src-$latestTag.tar.gz 
