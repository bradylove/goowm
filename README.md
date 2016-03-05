Goo Window Manager
======================================

GooWM is in very early days of development and should not be taken seriously, yet.

## Testing

For testing you can run GooWM using [Xephyr](https://www.freedesktop.org/wiki/Software/Xephyr/)

```
$ Xephyr :1 -ac -screen 1024x748 &
```

Then run GooWM.

```
$ DISPLAY=:1 go run main.go
```

Then start an application on display `:1`

```
$ DISPLAY=:1 xterm
```
