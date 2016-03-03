
# cam

Cam is a simple IP camera proxy: it polls single-frame JPEG output from one or
more IP cameras, and serves the frames out over HTTP in JPEG and MJPEG format.

*This was a quick and dirty solution developed for a friend.*  More serious
applications may benefit from [ffmpeg](http://ffmpeg.org/),
[gstreamer](http://gstreamer.freedesktop.org/), or a dedicated IP camera
package.


## Build

Clone the source:

    $ git clone https://github.com/lnmx/cam

    $ cd cam


Build the `cam` executable:

    $ export GOPATH=`pwd`

    # Windows: set GOPATH=%CD%

    $ go build cam


## Configure

Edit `config.json` to set up the listening port and source cameras, ex:  

    {
        "server": { 
            "addr": "0.0.0.0:8080" 
        },
        "sources": [
            { 
                "name": "camera1", 
                "url": "http://camera1/jpeg.cgi" 
                "refresh": 1.5
            },
            { 
                "name": "camera2", 
                "url": "http://camera2/jpeg.cgi", 
                "user": "test", 
                "pass": "secret",
                "refresh": 1
            }
        ]
    }

Include `user` and `pass` if the camera requires HTTP basic authentication.  

The `refresh` parameter sets the number of seconds between frames.


## Run

Run `cam`:

    $ ./cam

    2016/02/02 08:08:08 connecting camera2 to http://camera2/jpeg.cgi
    2016/02/02 08:08:08 connecting camera1 to http://camera1/jpeg.cgi
    2016/02/02 08:08:08 http server listening on 0.0.0.0:8080
    2016/02/02 08:08:08 connected to camera1
    2016/02/02 08:08:08 connected to camera2

    ...

Output URLs are based on the `sources/name` attribute in `config.json`.  Each
camera has a JPEG and MJPEG output.  For our example configuration:

  * JPEG: `http://cam-address:8080/camera1.jpeg`
  * MJPEG: `http://cam-address:8080/camera1.mjpeg`



