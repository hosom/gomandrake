# gomandrake

Mandrake is a file analysis framework. It is similar to Lockheed Martin's Laikaboss or Emerson's fsf. 

Mandrake monitors a directory for new files to be written, utilizing inotify, and then pumps those files through analyzers to perform static and unattended dynamic analysis. Mandrake strives to be both easy to set up and easy to manage.

## Installation

For now, Mandrake needs to be built from source. It is a short term goal to get a binary release for Ubuntu posted as a release.

```
# Installation of dependencies
sudo apt-get install golang libmagic-dev yara python-yara git

# Check out the code from github
git clone https://github.com/hosom/gomandrake

# Set a GOPATH for go to store the deps in
mkdir ~/MandrakeBuild
export GOPATH=~/MandrakeBuild

# Have go resolve required dependencies
go get github.com/hosom/gomandrake

# Build Mandrake
cd gomandrake
go build main.go

# Installation of pymandrake python library for python based plugins
sudo pip install git+https://github.com/hosom/pymandrake

# Instructions on how to build from source here

```
