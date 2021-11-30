## Build application

FROM golang:1.16-alpine
WORKDIR /app

# Copy go module files to download all required modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./

#Bulid file as "auth" and run
RUN go build -o /auth
CMD ["/auth"]

# ## Deploy app using distroless image
# FROM gcr.io/distroless/base-debian10
# WORKDIR /

# COPY --from=build /auth /auth

# EXPOSE 5003

# USER nonroot:nonroot

# ENTRYPOINT ["/auth"]