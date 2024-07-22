#!/bin/bash

SETUP_DIRECTORY="/Library/UserSetup"

if [[ -d "$SETUP_DIRECTORY" ]]; then
    echo "Setup automation directory found - remediating"
    exit 1
else
    echo "No setup automation directory found - no action needed"
    exit 0
fi