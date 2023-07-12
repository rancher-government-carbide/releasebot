#!/bin/sh
# install git hooks

PWD=$(git rev-parse --show-toplevel)
if [ "$1" = "link" ]; then
    printf "Symlinking hooks...\n"
    find "$PWD"/hooks/ -maxdepth 1 -exec basename {} \; | sed 1d | while read -r file; do
        ln -svf ../../hooks/"$file" "$PWD"/.git/hooks/"$file"
    done
else
    printf "Copying hooks...\n"
    find "$PWD"/hooks/ -maxdepth 1 -exec basename {} \; | sed 1d | while read -r file; do
        cp "$PWD"/hooks/"$file" "$PWD"/.git/hooks/"$file"
    done
fi
printf "Done!\n"
