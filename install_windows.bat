@echo off
rem run this script as admin

if not exist HelloWorldGoOsService.exe (
    echo Build the HelloWorldGoOsService before installing by running "go build"
    goto :exit
)

sc create go-svc-HelloWorldGoOsService binpath= "%CD%\HelloWorldGoOsService.exe" start= auto DisplayName= "go-svc-HelloWorldGoOsService"
sc description go-svc-HelloWorldGoOsService "go-svc-HelloWorldGoOsService"
net start go-svc-HelloWorldGoOsService
sc query go-svc-HelloWorldGoOsService

echo Check HelloWorldGoOsService.log

:exit
