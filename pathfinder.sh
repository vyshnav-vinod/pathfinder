#!/bin/sh

dir=$(./main "$@")

case $? in
    0)
    cd "$dir"
    ;;
    33)
    echo "No such dir"
    ;;
    *)
    echo "Noooo"
    ;;
esac