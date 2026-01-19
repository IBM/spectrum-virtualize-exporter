FROM ubuntu:jammy AS pem
WORKDIR /root/


FROM ubuntu:jammy

ARG APP_USER=spectrum

# Use "make binary" to build the binary spectrum-virtualize-exporter
COPY spectrum-virtualize-exporter /opt/spectrumVirtualize/spectrum-virtualize-exporter
COPY spectrumVirtualize.yml /opt/spectrumVirtualize/spectrumVirtualize.yml

RUN groupadd -g 1000 -r $APP_USER \
    && useradd -u 1000 -r -g $APP_USER -d /home/$APP_USER -m -s /bin/bash $APP_USER \
    && chown -R 1000:1000 /opt/spectrumVirtualize

USER $APP_USER

# port of prometheus exporter endpoint 
EXPOSE 9119

ENTRYPOINT ["/opt/spectrumVirtualize/spectrum-virtualize-exporter"]
CMD ["--config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml"]
