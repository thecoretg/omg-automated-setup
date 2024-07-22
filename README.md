# omg-automated-setup
Command-line setup assistant for OMG Mac setups. Deployed via Kandji in script: "Mac Setup Automation: Deploy" and removed after completion with "Mac Setup Automation: Remove" (in "deploy" and "remove" folders in this repository, respectively)

## What does it do?
Once deployed via Kandji, omgsetup goes onto any Mac that is in the "Ready to Assign" blueprint. It can then be run via typing `sudo omgsetup` in terminal, which will auto-detect if there are any users assigned to the Mac and then will go through a guided setup for the technician to go through. 

### If there is a user assigned:

1. Check if the user already exists, and if a spare user exists.
2. Checks if the assigned user already exists on the Mac and if so, skips the creation of the user.
3. Checks if a spare user exists, and if so, offers to remove it.
4. Menu will ask if the user is a standard or dev user.
5. A user will be created based as standard or admin based on #4 - with the standard default password.
6. Device is moved to the standard or dev blueprint.
7. Gives a summary of what was done.

### If there is no user assigned:

1. Checks if spare already exists - if so, exits.
2. Asks if you want to continue with spare user creation.
3. Creates a spare user with the default spare password.
4. Gives a summary of what was done.
5. Kept in Ready to Assign Blueprint.