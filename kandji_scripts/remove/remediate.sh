#!/bin/bash

SETUP_DIRECTORY="/Library/UserSetup"

rm -rf "$SETUP_DIRECTORY"

if [[ -d "$SETUP_DIRECTORY" ]]; then
    echo "Setup automation directory still present. Remediation needed!"
    exit 1
else
    echo "Setup automation directory successfully erased."
    exit 0
fi