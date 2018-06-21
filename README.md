# Go Report Card
This is a fork from `gojp/goreportcard`.

Currently it's rewritten as 2 separate cli commands.
- One to scan a local path repo and save report + assets into repo's location
- Another as part of docker image that allows using the image to scan a private remote repo. 


### Installation
##### To install local cmd
```
make install-local
```
This will build go binary and install linter dependencies.
Run binary from root with argument path to local repo:
```
./goreportcard ~/go/src/github.com/gojp/reportcard
```
Binary needs `assets/` and `templates/` folders to be on path from where it is executed. (hence the need for root)
this is the problem that needs solving with local scanning)

##### To install docker ci private repo cmd
```
make docker-ci
```
`WARNING`: This copies git ssh key from `~/.ssh/id_rsa` and adds it to the built image.
Then the image can be run with a remote public or private (if key was present/valid) repo argument, html with results goes to stdout, therefore can be appended with optional `> report.html`:
```
docker run meetcircle/goreportcard github.com/meetcircle/fakedbd > report.html
``` 
This can already be plugged into a CI pipeline, but has couple issues:   
1. The limitation of not being able to checkout any branch other than default  
		- Can be solved either by switching from `golang.org/x/tools/go/vcs` to a different one, say `go git`, that can checkout code by branch, as it seems.  
		- Or by mounting a volume of checked out branch by jenkins and providing a path to it as arg 
		(if this is properly possible - then the need for `local cmd` would be eliminated altogether)
		(ie: `docker run -v path/to/checkedout/repo:path/to/repo meetcircle/goreportcard path/to/repo > report.html`) - I have tried this approach, but for some reason linters werent scanning the files properly, report would always show correct number of files scanned but always 0 errors found. Have no time now to investigate this further, but this should be a working possibility. 
2. Only html is sent to stdout.   
		- Since `assets/` are already added to docker image. It should be possible to compress html and assets/ into a tar, spit it to stdout and collect via `> report.tar` (if my brief research on `tar + stdout` was correct)

### Development

Main files for both commands are in `cmd/goreportcard/`  
Both commands use a ReportHandler declared in `handlers/report.go`  
ReportHandlers use checking functions within `handlers/checks.go`  
###### Note
A lot of original code that supports web-app is still within the repo, perhaps it should be wiped altogether.

### License

The code is licensed under the permissive Apache v2.0 licence. This means you can do what you like with the software, as long as you include the required notices. [Read this](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) for a summary.

### Notes

We don't support this on Windows since we have no way to test it on Windows.
