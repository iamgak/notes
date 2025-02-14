# Base image: specify the environment
FROM golang:latest

# Set the working directory
# WORKDIR /app

# Copy the rest of the files
# COPY . /app

RUN mkdir /app

##set main directory
WORKDIR /app

##copy the whole file to app
ADD . /app

## get all the packages
RUN go mod download

##create executeable
RUN go build -o /app/main .

##run executeable
CMD ["/app/main"]