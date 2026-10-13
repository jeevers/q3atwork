FROM debian:bookworm-slim

COPY ioquake3 /usr/local/ioquake3
EXPOSE 27960/udp

#RUN adduser -D quakeserver
RUN useradd -m quakeserver
RUN chown -R quakeserver:quakeserver /usr/local/ioquake3
USER quakeserver

ENTRYPOINT ["/bin/sh"]
CMD ["/usr/local/ioquake3/start_server.sh"]
