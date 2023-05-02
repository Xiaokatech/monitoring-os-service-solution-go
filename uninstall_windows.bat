@echo off
rem run this script as admin

net stop go-svc-HelloWorldGoOsService
sc delete go-svc-HelloWorldGoOsService
