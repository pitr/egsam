FROM alpine:latest

ADD build/linux/egsam13 /
ADD egsam-prod.crt /egsam.crt
ADD egsam-prod.key /egsam.key

CMD ["/egsam13"]