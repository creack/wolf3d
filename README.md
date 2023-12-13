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

## Ack

Largely based on:
- https://lodev.org/cgtutor/raycasting.html#Untextured_Raycaster_
- https://github.com/faiface/pixel-examples/tree/704acac0e5f6fc19b27d5772033d77fc58cb7d59/community/raycaster
