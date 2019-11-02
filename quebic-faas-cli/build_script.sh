dist=$1
version=$2

# linux
linux_dist=$dist/linux
linux_file=$linux_dist/quebic-v$version-linux
mkdir $linux_dist
GOOS=linux GOARCH=amd64 go build -o $linux_file
echo "successfully build for linux"

# mac
macos_dist=$dist/macos
macos_file=$macos_dist/quebic-v$version-macos
mkdir $macos_dist
GOOS=darwin GOARCH=amd64 go build -o $macos_file
echo "successfully build for macos"

# windows
windows_dist=$dist/windows
windows_file=$windows_dist/quebic-v$version-windows
mkdir $windows_dist
GOOS=windows GOARCH=amd64 go build -o $windows_file
echo "successfully build for windows"

echo "build completed"

