# next2solve
Next problems to solve:
selection of [UVa](http://uva.onlinejudge.org/) problems
provided by the [uHunt](http://uhunt.felix-halim.net/) website.

## Try it! :smile:

https://s106.net/next2solve


## Installation and Usage

    mkdir src
    export GOPATH=$(pwd)
    cd src/
    git clone https://github.com/maru/next2solve
    cd next2solve/
    go build

If you want to use the original uHunt server, you can just run
(default listening port is 8002):

    ./next2solve -port 8002

You can also use a local testing server to provide the uHunt API responses
(see files in `testing` directory):

Terminal 1:

    cd testing/
    python -m SimpleHTTPServer 8080

Terminal 2:

    ./next2solve -api http://localhost:8080

Open http://localhost:8002/ in your browser.

## Testing

  go test next2solve next2solve/uhunt next2solve/problems
