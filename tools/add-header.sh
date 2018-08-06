#!/usr/bin/env bash

# Always points to the directory of this script.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..

# loop over all go files
echo "removing the old headers"
for i in $(find . -name \*.go); do
    python ./tools/replace-header.py $i > $i.tmp
    rm $i
done

# add new header to all tmp files and remove tmp
echo "add new header to file"
for itmp in $(find . -name \*.go.tmp); do
    cat $DIR/header.txt $itmp > $(echo $itmp | sed 's/.tmp//')
    rm $itmp
done

# store the old header
echo "store the old header"
cat $DIR/header-old.txt > $DIR/header.txt

# run go fmt
echo "run go fmt"
go fmt ./...