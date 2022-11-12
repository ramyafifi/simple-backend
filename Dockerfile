# Use AWS's Public ECR golang image
FROM public.ecr.aws/bitnami/golang:latest
ENV GO111MODULE=auto

# Set incoming environment variable to the workdir
COPY . .

#Specify output path to run application
RUN go build -o /bin/ab3

#Run application
CMD ["/bin/ab3"]