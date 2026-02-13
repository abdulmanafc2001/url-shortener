# Use distroless as minimal base image to package the manager binary
FROM gcr.io/distroless/static:nonroot
WORKDIR /

# Copy build file from local to image
COPY bin/manager /manager

EXPOSE 8080

CMD [ "/manager" ]