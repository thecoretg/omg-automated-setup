#!/bin/bash

# Create a new user on the mac with the given username and password, and add to admin if specified

# Check if the script is being run as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

# Check if the username and password are provided
if [ "$#" -ne 3 ]; then
  echo "Incorrect amount of arguments passed - usage: $0 <username> <password> <standard|admin>"
  exit 1
fi

# Check if the user already exists
if id "$1" &>/dev/null; then
  echo "User $1 already exists - please delete and then run again"
  exit 1
fi

# Create the user and make an admin if specified
if [ "$3" == "admin" ]; then
  dscl . -create /Users/"$1"
  dscl . -create /Users/"$1" UserShell /bin/bash
  dscl . -create /Users/"$1" RealName "$1"
  dscl . -create /Users/"$1" UniqueID 1001
  dscl . -create /Users/"$1" PrimaryGroupID 80
  dscl . -create /Users/"$1" NFSHomeDirectory /Users/"$1"
  dscl . -passwd /Users/"$1" "$2"
  dscl . -append /Groups/admin GroupMembership "$1"
else
  dscl . -create /Users/"$1"
  dscl . -create /Users/"$1" UserShell /bin/bash
  dscl . -create /Users/"$1" RealName "$1"
  dscl . -create /Users/"$1" UniqueID 1001
  dscl . -create /Users/"$1" PrimaryGroupID 20
  dscl . -create /Users/"$1" NFSHomeDirectory /Users/"$1"
  dscl . -passwd /Users/"$1" "$2"
fi

# shellcheck disable=SC2181
if [ "$?" -eq 0 ]; then
  echo "User $1 created successfully"
  exit 0
else
  echo "Failed to create user $1"
  exit 1
fi
