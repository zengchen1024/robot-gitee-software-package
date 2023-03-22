#!/bin/sh

set -euo pipefail

init() {
    if [ -d $repo ]; then
       return
    fi

    git clone --depth=1 $repo_url
    cd $repo

    git config user.name $git_user
    git config user.email $git_email
    git config --global pack.threads 1

    git remote add upstream ${upstream}
}

new_branch() {
    cd  $repo

    git checkout -- .
    git clean -fd

    git checkout master

    git fetch upstream master
    git rebase upstream/master

    git checkout -b $branch_name
}

commit() {
    cd $repo

    git add .

    git commit -m 'apply new package'

    git push origin $branch_name

    git checkout master
}

cmd=$1
git_user=$2
git_password=$3
git_email=$4
branch_name=$5

repo=community
upstream=https://gitee.com/openeuler/${repo}.git
repo_url=https://${git_user}:${git_password}@gitee.com/${git_user}/${repo}.git

case $cmd in
    "init")
        init
        ;;
    "new")
        new_branch
        ;;
    "commit")
        commit
        ;;
esac