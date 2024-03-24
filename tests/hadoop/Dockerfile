FROM ubuntu:latest

ARG HADOOP_VERSION=3.3.6

RUN apt-get update && apt-get install -y supervisor openjdk-11-jre krb5-kdc krb5-admin-server

ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64

ADD krb5.conf /etc/krb5.conf
RUN kdb5_util create -P password -r EXAMPLE.COM -s
RUN kadmin.local addprinc -pw password admin/admin@EXAMPLE.COM

ADD https://dlcdn.apache.org/hadoop/common/hadoop-${HADOOP_VERSION}/hadoop-${HADOOP_VERSION}.tar.gz /tmp/
RUN mkdir -p /var/lib/hadoop
RUN tar --strip-components 1 -xzf /tmp/hadoop-${HADOOP_VERSION}.tar.gz -C /var/lib/hadoop

RUN kadmin.local addprinc -pw password hadoop/localhost@EXAMPLE.COM
RUN kadmin.local addprinc -pw password HTTP/localhost@EXAMPLE.COM
RUN kadmin.local ktadd -norandkey -k /var/lib/hadoop/etc/hadoop/keytab hadoop/localhost@EXAMPLE.COM HTTP/localhost@EXAMPLE.COM

RUN mkdir -p /var/lib/hadoop-data
RUN mkdir -p /var/lib/hadoop-data/namenode
RUN mkdir -p /var/lib/hadoop-data/datanode
ADD core-site.xml /var/lib/hadoop/etc/hadoop/
ADD hdfs-site.xml /var/lib/hadoop/etc/hadoop/
RUN /var/lib/hadoop/bin/hdfs namenode -format

RUN mkdir -p /var/log/supervisor
COPY supervisord.conf /etc/supervisord.conf

EXPOSE 9870 9864 8088 8042 88
CMD ["/usr/bin/supervisord"]
