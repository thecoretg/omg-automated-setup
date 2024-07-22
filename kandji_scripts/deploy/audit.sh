#!/bin/bash

SETUP_BINARY="/Library/UserSetup/omgsetup"
SETUP_SYMLINK="/usr/local/bin/omgsetup"
SETUP_CONFIG="/Library/UserSetup/config.json"

if [[ ! -f "$SETUP_BINARY" || ! -f "$SETUP_CONFIG" || ! -f "$SETUP_SYMLINK" ]]; then
    echo "At least one setup file is missing - moving to remediation"
    exit 1
else
    echo "All needed files present - no action needed"
    exit 0
fi