#!/usr/bin/env bash

set -euo pipefail

scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

if [ $# != 1 ]; then
  echo Usage: make-image.sh IMAGE
  echo
  echo Image is one of:
    echo "  • all (build all images)"
  for d  in `ls $scriptdir/images`; do
    echo "  • $d"
  done
fi

images=$1
if [ "$images" == "all" ]; then
  images=""
  for d  in `ls $scriptdir/images`; do
    images+=" $d"
  done
fi

for image in $images; do
  (cd $scriptdir/images/$image && docker build -t quay.io/footloose/$image .)
done
