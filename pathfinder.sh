#!/bin/sh

# pfexecpath should be equal to the path of the pathfinder executable
dir=$(pfexecpath "$@")


case $? in
    0)
    cd "$dir"
    ;;
    1)
    echo "pf: Folder not found : $dir"
    ;;
    *)
    ;;
esac