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

if [[ $src_rpm_url == *"gitee.com"* ]]
then
  /opt/app/download $src_rpm_url ${ROBOT_TOKEN}
else
  curl -LO $src_rpm_url
fi

rpm2cpio *.rpm | cpio -div
rm -rf *.rpm

git add .

git commit -m 'add spec and rpm'

git push

cd ..

rm -rf $repo