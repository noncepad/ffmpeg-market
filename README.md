# FFmpeg-market

FFmpeg-Converter is server built to convert Blender files (.blend) into various video formats. This project utilizes a wrapper around FFMpeg in Go, enabling an efficient conversion process.

For details visit [noncepad.com](https://noncepad.com/blog/ffmpeg/)

# Features

* gRPC Streaming: Implements gRPC to stream files between client and server, ensuring efficient communication.

* Worker Pool: Utilizes a worker pool to handle conversion jobs concurrently, optimizing resource usage.

* Client-Server Architecture: Facilitates client-server communication, where clients upload .blend files along with specifications for desired output formats and output directory.

* Flexible Output: Allows clients to specify the desired output formats and directory for the converted files.

# Usage

To utilize the FFmpeg-Converter API, ensure that you have Blender installed and accessible in your environment. Clients can make requests by uploading .blend files along with the extension list for desired output formats and the directory to store the converted files.

# Build
Server:
```bash
go build -o ./ffmpeg-server github.com/noncepad/ffmpeg-market/cmd/server
```
Client:
```bash
go build -o ./ffmpeg-client github.com/noncepad/ffmpeg-market/cmd/client
```
# Run

Server:
```bash
./ffmpeg-server <myworkingdirectory> <listenurl for grpc> <MaxJobs>
```
Client:
```bash
./ffmpeg-client <serverURL> <filepathIn> <DirectoryOut> <extesion1> <extesion2> .... <extesionN>
```

## Example Arguments

Server:
```bash
./ffmpeg-server . localhost:8080 8 
```
Client:
```bash
./ffmpeg-client localhose :8080 ./video.blend . gif mp4
```