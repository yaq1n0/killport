@echo off
setlocal enabledelayedexpansion

REM KillPort Auto-Installer for Windows
REM This script downloads and installs killport automatically

echo.
echo 🔫 KillPort Auto-Installer
echo ==========================

REM Check if we're running as administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo ❌ This script requires administrator privileges.
    echo Please run as administrator and try again.
    pause
    exit /b 1
)

set BINARY_NAME=killport-windows-amd64.exe
set DOWNLOAD_URL=https://github.com/tarantino19/killport/releases/latest/download/%BINARY_NAME%
set INSTALL_DIR=C:\Windows\System32
set INSTALL_PATH=%INSTALL_DIR%\killport.exe
set TEMP_PATH=%TEMP%\killport.exe

echo 📋 Detected: Windows AMD64
echo 📥 Downloading: %BINARY_NAME%

REM Download using PowerShell
powershell -Command "& {Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%TEMP_PATH%'}"

if not exist "%TEMP_PATH%" (
    echo ❌ Download failed. Please check your internet connection.
    pause
    exit /b 1
)

echo ✅ Downloaded successfully
echo 📁 Installing to %INSTALL_PATH%

REM Copy to system directory
copy "%TEMP_PATH%" "%INSTALL_PATH%" >nul

if %errorLevel% neq 0 (
    echo ❌ Installation failed. Please check permissions.
    pause
    exit /b 1
)

REM Clean up temp file
del "%TEMP_PATH%" >nul 2>&1

echo ✅ killport installed successfully!
echo.
echo 🎯 Try it out:
echo    killport list
echo    killport 3000
echo.
echo 📖 Need help? Visit: https://github.com/tarantino19/killport
echo.
pause