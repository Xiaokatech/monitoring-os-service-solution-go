#!/bin/bash

sudo cp HelloWorldGoOsService.service /etc/systemd/system/

sudo systemctl enable HelloWorldGoOsService
sudo systemctl start HelloWorldGoOsService

sudo systemctl status HelloWorldGoOsService