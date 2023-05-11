@echo off
rem run this script as admin

if not exist HelloWorldGoOsServiceApp.exe (
    echo Build the HelloWorldGoOsServiceApp before installing by running "go build"
    goto :exit
)

sc create go-svc-HelloWorldGoOsService binpath= "C:\go\agentOsService\HelloWorldGoOsServiceApp.exe" start= auto DisplayName= "go-svc-HelloWorldGoOsService"
sc description go-svc-HelloWorldGoOsService "go-svc-HelloWorldGoOsService"
net start go-svc-HelloWorldGoOsService
sc query go-svc-HelloWorldGoOsService

echo Check HelloWorldGoOsServiceApp.log

:exit
