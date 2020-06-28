FROM alpine:latest

ADD build/linux/egsam /
ADD static /
ADD egsam.crt /
ADD egsam.key /

CMD ["/egsam"]