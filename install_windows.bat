@echo off
rem run this script as admin

if not exist ansysCSPAgentManagerServiceApp.exe (
    echo Build the ansysCSPAgentManagerServiceApp before installing by running "go build"
    goto :exit
)

sc create go-svc-AnsysCSPAgentManagerService binpath= "C:\go\agentOsService\ansysCSPAgentManagerServiceApp.exe" start= auto DisplayName= "go-svc-AnsysCSPAgentManagerService"
sc description go-svc-AnsysCSPAgentManagerService "go-svc-AnsysCSPAgentManagerService"
net start go-svc-AnsysCSPAgentManagerService
sc query go-svc-AnsysCSPAgentManagerService

echo Check ansysCSPAgentManagerServiceApp.log

:exit
