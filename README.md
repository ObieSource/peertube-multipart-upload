# peertube-multipart-upload
Upload a video to a [Peertube](https://joinpeertube.org/) instance via a multipart upload using the REST API.

The REST API does offer an API endpoint to upload a video file in a single request, and this is pretty straightforward to implement (in a [bash script](https://gist.github.com/FiskFan1999/66daa3063f63418cb6957123d7f8955d), for example). However, there are many reasons that a single-request upload may not be ideal for some users. The upload maximum allowed file-size may impact the size of video that can be uploaded in this way. For those with slow internet that is prone to disconnections, having to spend a long time uploading a file just for it to lose connection and have to start from the beginning may be unacceptable.

For these and other reasons, uploading a file in multiple parts is preferable. Files can be uploaded in parts as small as a few megabytes each, and if a part fails that part itself can be tried again without losing too much time. Multipart uploading can also support concurrency, depending on the implementation. Indeed, the Peertube browser application executes a multi-part upload when the user uploads the file.

This script aims to simplify the process of multipart-uploads to be as easy as uploading the file via the browser application or via a single-part upload. It uses environmental variables to declare the required parameters, such as the username and password, and information about the video.

# Installation
- For production use, a compiled binary of the [latest stable release](https://github.com/ObieSource/peertube-multipart-upload/releases/latest) is available.
- peertube-multipart-upload can be installed via the go compiler
```bash
go install github.com/ObieSource/peertube-multipart-upload@v0.0.2
```

- It can also be compiled from source in the following way:
  - Install the go compiler from [here](https://go.dev/dl/) or by another method.
  - Clone the peertube-multipart-upload repository
  - Change working directory to the repository via and install the binary application via `go install`. This will install the necessary go modules.
```bash
git clone https://github.com/ObieSource/peertube-multipart-upload.git
cd peertube-multipart-upload/
go install
```

# Usage
Refer to [help.txt](https://raw.githubusercontent.com/ObieSource/peertube-multipart-upload/master/help.txt) or type `peertube-multipart-upload help`.

# Reference
Refer to the [Peertube API reference](https://docs.joinpeertube.org/api-rest-reference.html#operation/uploadResumableInit) for the documentation of the API calls used internally.

- [Initialize](https://docs.joinpeertube.org/api-rest-reference.html#operation/uploadResumableInit)
- [Upload](https://docs.joinpeertube.org/api-rest-reference.html#operation/uploadResumable)
- [Cancel](https://docs.joinpeertube.org/api-rest-reference.html#operation/uploadResumableCancel) (Likely won't be implemented, as multipart uploads are automatically canceled after a certain length of time if they are abandoned)

Note that a `/api/v1/videos/upload-resumable` call will automatically finish the upload when the final part is called.

# Chat
Please feel free to discuss this project on irc.liberta.casa#fiskfan1999 (TLS :6697) [(browser client)](https://liberta.casa/kiwi/#fiskfan1999).
