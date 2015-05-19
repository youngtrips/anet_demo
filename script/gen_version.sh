#/bin/bash

git rev-list HEAD | sort > config.git-hash
branch=`git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3`
LOCALVER=`wc -l config.git-hash | awk '{print $1}'`
if [ $LOCALVER \> 1 ] ; then
    VER=`git rev-list $branch | sort | join config.git-hash - | wc -l | awk '{print $1}'`
    if [ $VER != $LOCALVER ] ; then
        VER="$VER+$(($LOCALVER-$VER))"
    fi
    if git status | grep -q "modified:" ; then
        VER="${VER}M"
    fi
    VER="$VER $(git rev-list HEAD -n 1 | cut -c 1-7)"
    GIT_VERSION=r$VER
else
    GIT_VERSION=
    VER="x"
fi
rm -f config.git-hash

VERSION="$GIT_VERSION $branch"
echo $VERSION
