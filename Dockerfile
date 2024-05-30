FROM ubuntu:jammy AS pem
RUN apt update && apt install curl -y
WORKDIR /root/
# The links of the IBM root CA and intermediate certs are from https://daymvs1.pok.ibm.com/ibmca/certificates.do;
# for other well known CAs, just install the ca-certificates package.
RUN curl https://daymvs1.pok.ibm.com/ibmca/downloadCarootCert.do?file=carootcert.der --output carootcert.der --cipher 'DEFAULT:!DH' && \
    curl https://daymvs1.pok.ibm.com/ibmca/downloadCarootCert.do?file=caintermediatecert.der --output caintermediatecert.der --cipher 'DEFAULT:!DH' && \
    openssl x509 -inform der -in carootcert.der -out 01-carootcert.pem && \
    openssl x509 -inform der -in caintermediatecert.der -out 02-caintermediatecert.pem


FROM ubuntu:jammy

ARG APP_USER=spectrum

# Use "make binary" to build the binary spectrum-virtualize-exporter
COPY spectrum-virtualize-exporter /opt/spectrumVirtualize/spectrum-virtualize-exporter
COPY spectrumVirtualize.yml /opt/spectrumVirtualize/spectrumVirtualize.yml
COPY --from=pem /root/*.pem /usr/local/share/ca-certificates/
# https://github.com/golang/go/blob/master/src/crypto/x509/root_linux.go
RUN mkdir -p /etc/ssl/certs \
    && cat /usr/local/share/ca-certificates/*.pem >> /etc/ssl/certs/ca-certificates.crt \
    && groupadd -g 1000 -r $APP_USER \
    && useradd -u 1000 -r -g $APP_USER -d /home/$APP_USER -m -s /bin/bash $APP_USER \
    && chown -R 1000:1000 /opt/spectrumVirtualize

USER $APP_USER
EXPOSE 9119
ENTRYPOINT ["/opt/spectrumVirtualize/spectrum-virtualize-exporter"]
CMD ["--config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml"]
