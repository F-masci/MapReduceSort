:: Close system tasks
@echo off
setlocal enabledelayedexpansion

:: List of ports to check
set ports=45978 45979 45980 45981 45989 45990

:: Loop through each port
for %%p in (%ports%) do (
    echo Checking port %%p...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%%p"') do (
        set pid=%%a
        echo Found process ID !pid! using port %%p
        taskkill /PID !pid! /F
        echo Process !pid! has been terminated.
    )
)

endlocal