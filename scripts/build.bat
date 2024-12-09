@echo off

:: Define the target operating systems and architectures with file extensions
set GOOS=windows
set GOARCH=amd64
set FILE_EXTENSION=zip

:: Set the output directory
set OUTPUT_DIR=..\builds

:: Create the output directory if it doesn't exist
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

:: Build file name
set FILE_NAME=spoty-%GOOS%-%GOARCH%
set OUTPUT_FILE=%OUTPUT_DIR%\%FILE_NAME%

echo Building for %GOOS%/%GOARCH%...

:: Compile Go application
set GO_ENV= 
call set GOOS=%GOOS%
call set GOARCH=%GOARCH%
go build -o "%OUTPUT_FILE%.exe" ..\cmd\web

:: Compress to ZIP if the build succeeded
if exist "%OUTPUT_FILE%.exe" (
    powershell Compress-Archive -Path "%OUTPUT_FILE%.exe" -DestinationPath "%OUTPUT_FILE%.zip"
    del "%OUTPUT_FILE%.exe"
    echo Build completed: %OUTPUT_FILE%.zip
) else (
    echo Build failed. Executable file was not created.
)

pause
