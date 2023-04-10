#!/bin/sh

set -euo pipefail

repo_url=$1
repo=$2
user=$3
email=$4
spec_url=$5
src_rpm_url=$6

if [ -d $repo ]; then
    rm -rf $repo
fi

git clone $repo_url

cd $repo

git config user.username $user
git config user.email $email

curl -LO $spec_url

curl -LO $src_rpm_url
rpm2cpio *.rpm | cpio -div
rm -rf *.rpm

git add .

git commit -m 'add spec and rpm'

git push

cd ..

rm -rf $repo