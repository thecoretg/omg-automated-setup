#!/bin/bash

githubLink="https://github.com/tctgdanny/omg-automated-setup/raw/main/cmd/main/omgsetup"
directoryPath="/Library/UserSetup"
configPath="${directoryPath}/config.json"
binaryPath="${directoryPath}/omgsetup"

# Download latest binary from github
mkdir -p /Library/UserSetup/
sudo curl -L -o "$binaryPath" "$githubLink" && echo "Successfully downloaded package from github"

# Make binary executable
chmod +x "$binaryPath"
ln -s "$binaryPath" /usr/local/bin/omgsetup

# Create config JSON
cat <<EOF > "$configPath"
{
    "KANDJI_API_TOKEN": "API_TOKEN_HERE",
    "KANDJI_API_URL": "API_URL_HERE",
    "KANDJI_STANDARD_BLUEPRINT": "STANDARD_BLUEPRINT_ID_HERE",
    "KANDJI_DEV_BLUEPRINT": "DEV_BLUEPRINT_ID_HERE",
    "TEMP_PASSWORD": "TEMP_PASSWORD_HERE",
    "SPARE_PASSWORD": "SPARE_PASSWORD_HERE"
}
EOF

# Set permissions for config 
chown root "$configPath"
chmod 400 "$configPath"