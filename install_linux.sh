#!/bin/bash

sudo cp AnsysCSPAgentManagerService.service /etc/systemd/system/

sudo systemctl enable AnsysCSPAgentManagerService
sudo systemctl start AnsysCSPAgentManagerService

sudo systemctl status AnsysCSPAgentManagerService