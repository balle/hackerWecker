# hackerWecker
A tool to wake up a hacker in the morning by creating some kind of radio show (reading computer related rss / atom news, playing some music / podcast)

Run with `go run main/hackerWecker.go`

## Implementation

Although the tool runs on Linux, *BSD, macOS and Windows laptop or desktop I recommend implementing it on a [Rasperry Pi][https://www.raspberrypi.org/].
Note that you need to add a better sound card via USB as the integrated one is not capable in playing music without hurting ones ears.

To make sure that the sound is really played on the second sound card added via USB edit or create the file /etc/asound.conf with the following content:

`pcm.usb_snd_card {
    type hw
    card 1
    device 0
}

ctl.usb_snd_card {
    type hw
    card 1
    device 0
}

pcm.!default {
    type plug
    slave.pcm "usb_snd_card"
}`

To make the hackerWecker actually wake you up add a crontab line like the following. It will wake you up every workday on 6:30.

> # crontab -e

`30 6 * * 1-5 cd /home/pi/hackerWecker; /usr/bin/go run main/hackerWecker.go`
