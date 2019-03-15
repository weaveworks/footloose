#!/usr/bin/env bash

set -euo pipefail

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
org=quay.io/footloose

if [ $# -lt 2 ]; then
  echo "Usage: make-image.sh VERB IMAGE [ARGS]"
  echo
  echo VERB is one of:
  echo "  • build"
  echo "  • tag"
  echo "  • push"
  echo
  echo IMAGE is one of:
  echo "  • all (build all images)"
  for d  in `ls $scriptdir/images`; do
    echo "  • $d"
  done
  echo
  echo "ARGS are VERB-specific optional arguments"
  echo "  • tag requires a version"
  echo "  • push takes an optional version, defaults to 'latest' "
  echo
  echo "Examples:"
  echo "  • $0 build all"
  echo "  • $0 tag all 0.2.0"
  echo "  • $0 push all 0.2.0"
  exit 1
fi

verb=$1
case $verb in
build|tag|push)
  ;;
*)
  echo "error: unknown verb '$verb'"
  exit 1
esac

images=$2
if [ "$images" == "all" ]; then
  images=""
  for d  in `ls $scriptdir/images`; do
    images+=" $d"
  done
fi

shift
shift

case $verb in

build)
  for image in $images; do
    echo "  • Building $org/$image"
    (cd $scriptdir/images/$image && docker build -t $org/$image .)
  done
  ;;

tag)
  if [ $# != 1 ]; then
    echo "error: usage tag IMAGE VERSION"
    exit 1
  fi
  version=$1
  for image in $images; do
    echo "  • Tagging $org/$image:$version"
    docker tag $org/$image:latest $org/$image:$version
  done
  ;;

push)
  version=latest
  [ $# == 1 ] && version=$1
  for image in $images; do
    echo "  • Pushing $org/$image:$version"
    docker push $org/$image:$version
  done
  ;;

esac
