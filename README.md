# Payment QR code generator

Standalone web server generating payment QR codes.

Supports [Short Payment Descriptor](https://en.wikipedia.org/wiki/Short_Payment_Descriptor) and
[EPC QR Code](https://en.wikipedia.org/wiki/EPC_QR_code) (SEPA) format.

## Demo app

https://go-qrpay.com/

## Instalation

    go get -u github.com/dundee/gdu/qrpay-web

## Usage

1. Start the QR code generator:

    $ ./qrpay-web
    2021/04/02 00:41:01 using embed mode
    2021/04/02 00:41:01 Listening on 127.0.0.1:8080

1. Open `127.0.0.1:8080` in your web broswer.
1. Fill in payment details and submit.
1. See the result QR code.

![Screenshot](./screenshot.png)
