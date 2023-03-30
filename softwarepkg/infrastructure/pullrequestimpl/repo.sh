#!/bin/sh

set -euo pipefail

init() {
    if [ -d $repo_dir ]; then
       return
    fi

    git clone --depth=1 $repo_url
    cd $repo_dir

    git config user.name $git_user
    git config user.email $git_email
    git config --global pack.threads 1

    git remote add upstream ${upstream}
}

new_branch() {
    cd  $repo_dir

    git checkout -- .
    git clean -fd

    git checkout master

    git fetch upstream master
    git rebase upstream/master

    git checkout -b $branch_name
}

modify() {
  cd $repo_dir

  echo "$sig_info_content" >> $sig_info_file

  dn=$(dirname $new_repo_file)
  if [ ! -d $dn ]; then
     mkdir $dn
  fi

  echo "$new_repo_content" > $new_repo_file
}

commit() {
    cd $repo_dir

    git add .

    git commit -m 'apply new package'

    git push origin $branch_name

    git checkout master
}

git_user=$1
git_token=$2
git_email=$3
branch_name=$4
org=$5
repo=$6
sig_info_file=$7
sig_info_content=$8
new_repo_file=$9
new_repo_content=${10}

upstream=https://gitee.com/${org}/${repo}.git
repo_url=https://${git_user}:${git_token}@gitee.com/${git_user}/${repo}.git
repo_dir=$(pwd)/${repo}

init

new_branch

modify

commit