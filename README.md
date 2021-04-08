convertallaudio
===

Scans a directory recursively for wav files and converts them to a specific output format using FFMPEG(via FFMPEG environment variable).

## Usage
```
convertallaudio.exe --input=<some directory> --format=[flag,ogg,mp3] --output=<target directory>
```
### Run & Build
```
go run convertallaudio.go
go build
```
