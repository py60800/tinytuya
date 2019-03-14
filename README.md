# TinyTuya

TinyTuya is a minimal application dedicated to the control of Tuya switches.

Tuya is a Chinese [Tuya](https://en.tuya.com/) company that claims to be a leader in IoT (Internet of Things). Actually many Wifi home devices use Tuya protocols. They are sold with different brand names.

Android applications like Smart Life control devices using protocol devised by Tuya.

I have developed this small application because I wanted to be able to control my devices without relying solely on my smartphone (i.e. use my notebook).

TinyTuya achieve this goal in a very simple manner : the application can be installed on a Raspberry PI (any Linux or Windows box will do the job provided it can be kept running).

It is dedicated to home use when security is not a concern (anyone who has access to your network will be able to control the devices).

## Prerequisites

### Get the encryption keys

The control of Tuya devices requires to get the encryption keys that allow communication with the devices. @Codetheweb has developed tools to get these keys. Check [@Codetheweb page](https://github.com/codetheweb/tuyapi/blob/master/docs/SETUP.md)

Actually, any proxy that has the ability to decipher HTTPS can do the job… (not that simple)

## Golang environment

TinyTuya is written in GO, so you need to install GO tools to build the application.

## Build the application

- get communication library:

`go get github.com/py60800/tuya`

- get TinyTuya

`go get github.com/py60800/tinytuya`

- Build the application

depending on where the source code is:

`go build tinytuya.go`

## Configure and run

### Configure
The configuration file is JSON formatted. Use the example configuration and set the parameters according to you environment :
- ids of the tuya devices (collected at the same time as the keys)
- key
- friendly name

### Run

Change to the directory that contains additional resources (static, tmpl)

`./tinytuya -c config.json`

Default port is 8080, another port can be set with ‘-p’ option


## Use

From any PC, smartphone or tablet, connect to your device http://<Ip or Name>:8080

## Finalize

Adapt to your needs : appearance can be changed by adapting the templates files.

Out of the scope of this memo : run it as a service

## Limits

As of now :

- It has been tested with switches branded by Neo

- only version 3.1 of the protocol is supported

## Acknowledgments

@Codetheweb for this tremendous reverse engineering job

