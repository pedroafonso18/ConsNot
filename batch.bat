@echo off
:loop
    echo Starting Go process...
    go run cmd/ConsNot/main.go | findstr /I "bad connection"
    echo Detected bad driver error! Restarting...
    timeout /t 2
goto loop