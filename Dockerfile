FROM ubuntu:22.04 AS pem
RUN apt update && apt install curl -y
WORKDIR /root/
# The links of the IBM root CA and intermediate certs are from https://daymvs1.pok.ibm.com/ibmca/certificates.do;
# for other well known CAs, just install the ca-certificates package.
RUN curl https://daymvs1.pok.ibm.com/ibmca/downloadCarootCert.do?file=carootcert.der --output carootcert.der --cipher 'DEFAULT:!DH' && \
    curl https://daymvs1.pok.ibm.com/ibmca/downloadCarootCert.do?file=caintermediatecert.der --output caintermediatecert.der --cipher 'DEFAULT:!DH' && \
    openssl x509 -inform der -in carootcert.der -out 01-carootcert.pem && \
    openssl x509 -inform der -in caintermediatecert.der -out 02-caintermediatecert.pem


FROM busybox:latest
COPY spectrum-virtualize-exporter /bin/spectrum-virtualize-exporter
COPY spectrumVirtualize.yml /etc/spectrumVirtualize/spectrumVirtualize.yml
COPY --from=pem /root/*.pem /usr/local/share/ca-certificates/
# https://github.com/golang/go/blob/master/src/crypto/x509/root_linux.go
RUN mkdir -p /etc/ssl/certs && \
    cat /usr/local/share/ca-certificates/*.pem >> /etc/ssl/certs/ca-certificates.crt
EXPOSE 9119
ENTRYPOINT ["/bin/spectrum-virtualize-exporter"]
CMD ["--config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml"]
