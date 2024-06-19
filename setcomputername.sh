#!/bin/bash

# Get assigned user from Kandji global variables plist, and computer name from system_profiler
ASSIGNED_USER="$(/usr/libexec/PlistBuddy -c 'print :FULL_NAME' /Library/Managed\ Preferences/io.kandji.globalvariables.plist)"
COMPUTER_NAME="$(/usr/sbin/system_profiler SPHardwareDataType | awk -F": " '/Model Name/ {print $2}')"

# If assigned user is empty or if computer name is empty, exit
if [ -z "$ASSIGNED_USER" ] || [ -z "$COMPUTER_NAME" ]; then
    echo "Assigned user or computer name is empty. Exiting..."
    exit 1
fi

# If computer name is Mac mini, change it to Mac Mini for consistency in capitalization
if [ "$COMPUTER_NAME" == "Mac mini" ]; then 
    COMPUTER_NAME="Mac Mini"
fi

newComputerName="$ASSIGNED_USER $COMPUTER_NAME"
# shellcheck disable=SC2001
newHostName=$(echo "$newComputerName" | sed 's/ /-/g')

# Set computer name, local host name, and host name
scutil --set ComputerName "${newComputerName}" 
scutil --set HostName "${newHostName}"
scutil --set LocalHostName "${newHostName}" 

currentComputerName="$(scutil --get ComputerName)"
currentHostName="$(scutil --get HostName)"
currentLocalHostName="$(scutil --get LocalHostName)"

echo "Computer name is now: $currentComputerName"
echo "Host name is now: $currentHostName"
echo "LocalHostName is now: $currentLocalHostName"

