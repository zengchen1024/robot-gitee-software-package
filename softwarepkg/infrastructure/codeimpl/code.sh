#!/bin/sh

set -euo pipefail

# repo_url is the url to push code, it contains username and token
repo_url=$1
# repo is the repo name of repo_url
repo=$2
user=$3
email=$4
# get code from the pr of ci repo, therefore it can guarantee that the code is CI checked
ci_repo_link=$5
ci_repo=$6
ci_pr_num=$7

if [ ! -d $ci_repo ]; then
  git clone -q $ci_repo_link
fi

cd $ci_repo

git checkout master

branch_name="pr$ci_pr_num"

git fetch origin "pull/$ci_pr_num/head:$branch_name"

git checkout $branch_name

cd ..

if [ -d $repo ]; then
    rm -rf $repo
fi

git clone $repo_url

cp $ci_repo/*.spec $repo
cp $ci_repo/*.rpm $repo

cd $repo

rpm2cpio *.rpm | cpio -div
rm -rf *.rpm

git config user.username $user
git config user.email $email
git add .
git commit -m 'add spec and rpm'
git push

cd ..

rm -rf $repo