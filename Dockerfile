FROM busybox:1.29

COPY dice /bin/dice

WORKDIR    /home
ENTRYPOINT "/bin/dice"
