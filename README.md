# Wolf3D - Simple Raycaster

## Live Demo

[https://creack.github.io/wolf3d](https://creack.github.io/wolf3d)

## Window mode

```sh
go run go.creack.net/wolf3d@latest
```

## WASM

### One liner

```sh
env -i HOME=${HOME} PATH=${PATH} go run github.com/hajimehoshi/wasmserve@latest go.creack.net/wolf3d@latest
```

### Details

Install `wasmer`:

```sh
go install github.com/hajimehoshi/wasmserve@latest
```

Clone this repo:

```sh
git clone https://github.com/creack/wolf3d
cd wolf3d
```

Run:

```sh
env -i HOME=${HOME} PATH=${PATH} wasmserve .
```

For development, `wasmer` exposes an endpoint to do live reload.

I recommend [reflex](https://github.com/cespare/reflex). 

Install:

```sh
go install github.com/cespare/reflex@latest
```

Then, with `wasmer` running:

```sh
reflex curl -v http://localhost:8080/_notify
```

## Controls

- up/down w/s:  Move up/down.
- right/left: Turn right/left.
- a/d: Strife right/left.

## Docker

A Dockerfile is provided to build and run the WASM version.

### Build

```sh
docker build -t wolf3d .
```

### Regular run

To run the image, make sure to have:
  - `--rm` to avoid pollution
  - `-it` so the app receives signals
  - `-p` to expose the port 8080

Any changes to the code will require to re-build the image.

```sh
docker run --rm -p 8080:8080 -it wolf3d wasmserve .
```

You can then access the WASM page at the Docker ip on port 8080. If in doubt about the IP, it is likely localhost.

### Development run 

For development, you can add `-v $(pwd):/app` to mount the local directory in the Docker container, the server will hot-reload when file changes.

```sh
docker run --rm -p 8080:8080 -it -v $(pwd):/app wolf3d
```

## Ack

Largely based on:
- https://lodev.org/cgtutor/raycasting.html#Untextured_Raycaster_
- https://github.com/faiface/pixel-examples/tree/704acac0e5f6fc19b27d5772033d77fc58cb7d59/community/raycaster
