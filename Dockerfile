FROM bitnami/minideb:latest
COPY ./minio-broker /minio-broker
ADD ./minio-0.1.3.tgz /minio-chart.tgz
CMD ["/minio-broker", "-logtostderr"]
