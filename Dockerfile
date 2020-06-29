FROM alpine:latest

ADD build/linux/egsam /
ADD static /static
ADD egsam-prod.crt /egsam.crt
ADD egsam-prod.key /egsam.key

CMD ["/egsam"]