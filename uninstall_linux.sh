# Copyright (C) 2023 ANSYS, Inc. Unauthorized use, distribution, or duplication is prohibited.

#!/bin/bash

sudo systemctl stop AnsysCSPAgentManagerService # Stop the running service
sudo systemctl disable AnsysCSPAgentManagerService # Disable the service so it doesn't start automatically on system boot

sudo rm /etc/systemd/system/AnsysCSPAgentManagerService.service

sudo systemctl daemon-reload # reload the systemd daemon to apply the changes
