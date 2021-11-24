#!/bin/sh

binName="odo3"
outputDir="dist/bin"

for platform in linux-amd64 darwin-amd64 darwin-arm64 windows-amd64; do
  echo "Cross compiling $platform and placing binary at $outputDir/$platform/"

  if [ ${platform%-*} = "windows" ]; then
    binName="$binName.exe"
  fi

    GOARCH=${platform#*-} GOOS=${platform%-*} go build -o $outputDir/$platform/$binName "${@}" main.go

  if [ ${platform%-*} = "windows" ]; then
    zip -j $outputDir/$platform-$binName.zip $outputDir/$platform/$binName
  else
    gzip -c $outputDir/$platform/$binName > $outputDir/$platform-$binName.gz
  fi

done
