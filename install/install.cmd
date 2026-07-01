@echo off
setlocal
if exist "%~dp0install.ps1" (
  powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1" %*
  exit /b %ERRORLEVEL%
)
set "QUALITYMD_INSTALLER=%TEMP%\qualitymd-install-%RANDOM%-%RANDOM%.ps1"
powershell -NoProfile -ExecutionPolicy Bypass -Command "Invoke-WebRequest https://getquality.md/install.ps1 -UseBasicParsing -OutFile '%QUALITYMD_INSTALLER%'"
if errorlevel 1 exit /b %ERRORLEVEL%
powershell -NoProfile -ExecutionPolicy Bypass -File "%QUALITYMD_INSTALLER%" %*
set "QUALITYMD_EXIT=%ERRORLEVEL%"
del "%QUALITYMD_INSTALLER%" >nul 2>nul
exit /b %QUALITYMD_EXIT%
