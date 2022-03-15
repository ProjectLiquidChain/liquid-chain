#!/bin/bash
root=`git rev-parse --show-toplevel`
green=`tput setaf 2`
reset=`tput sgr0`

while read line; do
    arr=(${line//=/ })
    filepath=${arr[0]}
    exports=${arr[1]}
    echo "$green$filepath$reset"
    cd $root/`dirname $filepath`
    vertex-cdt compile `basename $filepath` --export-function=$exports
done < $root/wasm-funcs.export