# For building project

## Environment Variables
You need to create a `.env` file in the root directory of the project. Use the following structure as an example:
```
APP_PORT=8080

SPOTIFY_CLIENT_ID=spotify-client-id
SPOTIFY_CLIENT_SECRET=spotify-client-secret
SPOTIFY_REDIRECT_URI=http://localhost:8080/spotify/callback
```

## Windows
```
cd .\scripts && build.bat
```

## Linux or MacOS
```
cd ./scripts && sh build.sh
```
