# goCrunchQueue

Golang Backend and Vue Frontend to put different Crunchyroll URL's into a Queue to get them Processesd
by [crunchy-cli](https://github.com/crunchy-labs/crunchy-cli).

## Deployment

### Docker

Use  [serverfrog/gocrunchqueue](https://hub.docker.com/r/serverfrog/gocrunchqueue) as already existing image with
[crunchy-cli](https://github.com/crunchy-labs/crunchy-cli), ffmpeg and goCrunchQueue with its UI.
Use the Enviroment Variable `ETP_RT` or `CREDENTIALS` to specify them
for [crunchy-cli (Section Login)](https://github.com/crunchy-labs/crunchy-cli#login). Create and configure the
config.yaml which will be supplied to the Application. For an Example
view  [here](https://github.com/Serverfrog/goCrunchQueue/blob/main/config/config.yaml).

Example Docker Run Command:

`docker run  --name goCrunchQueue -p 80:80 --env ETP_RT=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee -v ./local-media:/goCrunchQueue/media-destination -v ./local-logs:/goCrunchQueue/logs gocrunchyqueue`

This will start the Container where the Frontend will be presented on `http://localhost:80` and the Downloaded Media will be in `./local-media`.

### Config

| Value            | default             | -                                                           |
|------------------|---------------------|-------------------------------------------------------------|
| Debug            | true                | will set the Log to Debug Mode                              |
| Port             | 80                  | the HTTP Port where the application will listen to.         |
| MediaDestination | ./media-destination | the Location where crunchy-cli will download its content to |
| LogDestination   | ./logs              | where the Logfiles will be written to                       |

