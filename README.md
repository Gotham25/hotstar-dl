# hotstar-dl
Golang cli app to download hotstar videos

[![Build Status](https://travis-ci.org/Gotham25/hotstar-dl.svg?branch=master)](https://travis-ci.org/Gotham25/hotstar-dl) [![Coverage Status](https://coveralls.io/repos/github/Gotham25/hotstar-dl/badge.svg?branch=v1.0.0)](https://coveralls.io/github/Gotham25/hotstar-dl?branch=v1.0.0) [![Go Report Card](https://goreportcard.com/badge/github.com/Gotham25/hotstar-dl)](https://goreportcard.com/report/github.com/Gotham25/hotstar-dl)

Static builds are found under `Assets` section in `Releases` tab.

Use args such as -h or --help to view the usage of application.

#### Sample Demo
![hotstar-dl_v2 1 0_demo](https://user-images.githubusercontent.com/12382378/84975183-68a7b280-b142-11ea-9e78-7b721c58bff2.gif)

#### Steps to download video/audio
1. Download binaries from [here](https://github.com/Gotham25/hotstar-dl/releases)
2. Get a sample url from hotstar website
3. To list the available formats to download use the below command,
   
   hotstardl.exe -l \<URL\>
   
   where URL is the sample URL from step 1
4. Choose a format from the above list.
5. To download video/audio use the below command
   
   hotstardl.exe -f \<FORMAT\> \<URL\>
   
   where
   - URL&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;sample URL from step 1 and 
   - FORMAT&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;format choosen in step 4
