FROM alpine:latest

ADD build/linux/egsam /
ADD static /static
ADD egsam.crt /
ADD egsam.key /

CMD ["/egsam"]