@REM Copyright (C) 2023 ANSYS, Inc. Unauthorized use, distribution, or duplication is prohibited.

@echo off
rem run this script as admin

net stop go-svc-AnsysCSPAgentManagerService
sc delete go-svc-AnsysCSPAgentManagerService
