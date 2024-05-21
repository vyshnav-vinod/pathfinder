#!/bin/bash

exec_name="bin/pathfinder"
current_dir=$(pwd)
exec_path="$current_dir/$exec_name"

if [[ ! -f "$exec_path" ]]; then
  echo "Error: $exec_name not found in the current directory"
  exit 1
fi

script_content="#!/bin/sh

dir=\$($exec_path \"\$@\")

# EXIT CODES
#  0 - Success
#  1 - Folder not found
#  4 - Cache cleaned successfully
#  5 - Info returned
# -1 - Error

case \$? in
    0)
    cd \"\$dir\"
    ;;
    1)
    echo \"pf: Folder not found : \$dir\"
    ;;
    4)
    echo \"Cache cleaned\"
    ;;
    5)
    echo \$dir
    ;;
    *)
    ;;
esac
"

pathfinder_script_path="$current_dir/pathfinder.sh"
echo "$script_content" > "$pathfinder_script_path"

chmod +x "$pathfinder_script_path"

alias_command="alias pf='. $pathfinder_script_path'"

if ! grep -Fxq "$alias_command" ~/.bashrc; then
  echo "$alias_command" >> ~/.bashrc
  echo "Alias 'pf' added to .bashrc"
else
  echo "Alias 'pf' already exists in .bashrc"
fi

echo "pathfinder.sh created successfully"
