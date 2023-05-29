@echo off
rem run this script as admin

net stop go-svc-AnsysCSPAgentManagerService
sc delete go-svc-AnsysCSPAgentManagerService
