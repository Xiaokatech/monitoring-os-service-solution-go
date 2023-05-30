# Copyright (C) 2023 ANSYS, Inc. Unauthorized use, distribution, or duplication is prohibited.

#!/bin/bash

sudo cp AnsysCSPAgentManagerService.service /etc/systemd/system/

sudo systemctl enable AnsysCSPAgentManagerService
sudo systemctl start AnsysCSPAgentManagerService

sudo systemctl status AnsysCSPAgentManagerService