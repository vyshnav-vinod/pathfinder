#!/bin/sh

# pfexecpath should be equal to the path of the pathfinder executable
dir=$(pfexecpath "$@")

# EXIT CODES
#  0 - Success
#  1 - Folder not found
#  4 - Cache cleaned successfully
# -1 - Error

case $? in
    0)
    cd "$dir"
    ;;
    1)
    echo "pf: Folder not found : $dir"
    ;;
    4)
    echo "Cache cleaned"
    ;;
    *)
    ;;
esac