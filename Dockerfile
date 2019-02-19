
# FROM quay.io/prometheus/busybox:latest
FROM s390x/busybox:latest
COPY spectrum-virtualize-exporter /bin/spectrum-virtualize-exporter
COPY spectrumVirtualize.yml /etc/spectrumVirtualize/spectrumVirtualize.yml
EXPOSE 9119
ENTRYPOINT ["/bin/spectrum-virtualize-exporter"]
CMD ["--config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml"]