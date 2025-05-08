dev_image_name := "gotify-telegram-relay-dev"

build go_os="linux" go_arch="amd64":
    docker run --rm -it \
        -v {{justfile_dir()}}:/src \
        --workdir /src \
        {{dev_image_name}} /bin/bash -c 'export GOOS={{go_os}} && export GOARCH={{go_arch}} && go build -o /src/build/gotify-telegram-relay'

build-image:
    docker build -t {{dev_image_name}} .

shell: build-image
    docker run --rm -it \
        -v {{justfile_dir()}}:/src \
        --workdir /src \
        {{dev_image_name}} /bin/bash
