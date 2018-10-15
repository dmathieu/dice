FROM k8s.gcr.io/debian-base-amd64:0.3.2

COPY dice /home
WORKDIR /home
