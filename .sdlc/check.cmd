
@echo OFF
SETLOCAL ENABLEEXTENSIONS
SET "script_name=%~n0"
SET "script_path=%~0"
SET "script_dir=%~dp0"
rem # to avoid invalid directory name message calling %script_dir%\config.bat
cd %script_dir%
call config.bat
cd ..
set project_dir=%cd%

set module_name=%REPO_HOST%/%REPO_OWNER%/%REPO_NAME%
set exe_path=bin\%REPO_NAME%.exe

echo script_name   %script_name%
echo script_path   %script_path%
echo script_dir    %script_dir%
echo project_dir   %project_dir%
echo module_name   %module_name%
echo exe_path      %exe_path%

cd %project_dir%

go test -race ./...
go test -cover ./...

REM for /f %%x in ('dir /AD /B /S lib') do (
REM     echo --- %%x
REM     cd %%x
REM     go test -mod vendor -cover ./...
REM )
REM
REM cd %project_dir%
REM
REM IF EXIST %exe_path% DEL /F %exe_path%
REM
REM @echo ON
REM call go build -mod vendor -ldflags "-s -X %module_name%/lib/core.Version=%APP_VERSION% -X %module_name%/lib/core.BuildTime=%TIMESTAMP% -X %module_name%/lib/core.GitCommit=win-dev-commit" ^
REM   -o %exe_path% "%module_name%/cmd/%REPO_NAME%"
REM
REM @echo OFF
REM for /f %%x in ('dir /AD /B /S cmd') do (
REM     echo --- %%x
REM     cd %%x
REM     go test -mod vendor -cover ./...
REM )
