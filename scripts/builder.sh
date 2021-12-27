#!/bin/bash

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64")

empty=""

for dict in ./cmd/*; do
  if [ -d "${dict}" ]; then
    project="${dict//.\/cmd\/$empty}"
    path="${dict}/main.go"

    for platform in "${platforms[@]}"
    do
        platform_split=(${platform//\// })
        GOOS=${platform_split[0]}
        GOARCH=${platform_split[1]}
        output_name=$project'-'$GOOS'-'$GOARCH
        if [ $GOOS = "windows" ]; then
          output_name+='.exe'
        fi

        echo '[Building]' $output_name

        env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -a -o "./dist/${output_name}" $path
        if [ $? -ne 0 ]; then
            echo 'An error has occurred! Aborting the script execution...'
            exit 1
        fi
    done
  fi
done

echo 'DONE!'
