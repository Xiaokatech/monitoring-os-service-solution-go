#!/bin/bash

sudo systemctl stop HelloWorldGoOsService # Stop the running service
sudo systemctl disable HelloWorldGoOsService # Disable the service so it doesn't start automatically on system boot

sudo rm /etc/systemd/system/HelloWorldGoOsService.service

sudo systemctl daemon-reload # reload the systemd daemon to apply the changes
