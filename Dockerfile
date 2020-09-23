FROM openvax/neoantigen-vaccine-pipeline
USER user
COPY target/breakbio-openvax /home/user/breakbio-openvax
RUN sudo chmod +x /home/user/breakbio-openvax
CMD ["/home/user/breakbio-openvax", "serve"]