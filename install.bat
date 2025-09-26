@echo off
setlocal

echo Building killport for Windows...
go build -o bin\killport.exe main.go

if not exist "%USERPROFILE%\bin" (
    echo Creating %USERPROFILE%\bin directory...
    mkdir "%USERPROFILE%\bin"
)

echo Installing killport to %USERPROFILE%\bin...
copy bin\killport.exe "%USERPROFILE%\bin\" >nul

echo Adding %USERPROFILE%\bin to PATH if not already present...
echo %PATH% | find /i "%USERPROFILE%\bin" >nul
if errorlevel 1 (
    echo Current PATH does not contain %USERPROFILE%\bin
    echo Please add %USERPROFILE%\bin to your PATH environment variable manually.
    echo.
    echo Steps:
    echo 1. Open System Properties ^(Windows + Pause^)
    echo 2. Click "Advanced system settings"
    echo 3. Click "Environment Variables"
    echo 4. Under "User variables", find "Path" and click "Edit"
    echo 5. Click "New" and add: %USERPROFILE%\bin
    echo 6. Click "OK" to save changes
    echo 7. Restart your command prompt
) else (
    echo %USERPROFILE%\bin is already in PATH
)

echo.
echo killport installed successfully!
echo You can now use 'killport' from anywhere in your command prompt.
echo.
echo Usage examples:
echo   killport list          - List all active ports
echo   killport 3000          - Kill process on port 3000
echo   killport 3000 4000     - Kill processes on multiple ports
echo   killport all           - Kill all port processes ^(with confirmation^)
echo.
echo To uninstall, delete: %USERPROFILE%\bin\killport.exe

pause
