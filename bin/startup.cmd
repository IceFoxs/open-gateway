@echo off
set GO_ENV=test
set BASE_DIR=%~dp0
set CONF_PATH=%BASE_DIR%config\conf.yaml
set WEB_PATH=%BASE_DIR%web/
set COMMAND=%BASE_DIR%opengateway.exe
%COMMAND%