FROM openvax/neoantigen-vaccine-pipeline
COPY target/breakbio-openvax /opt/breakbio-openvax
RUN chmod +x /opt/breakbio-openvax
CMD ["/opt/breakbio-openvax", "serve"]