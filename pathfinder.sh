#!/bin/sh

# Find a way to get the path of the pathfinder executable
dir=$(./main "$@")


case $? in
    0)
    cd "$dir"
    ;;
    *)
    echo "$?"
    ;;
esac