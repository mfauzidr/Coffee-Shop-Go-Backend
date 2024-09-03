# Coffee Shop Backend with Go Language

Coffee Shop Web is a simple website for managing a coffee shop. This application makes it easier for users if they want to create a coffee shop business.

## Table of Contents

1. [About](#about)
   - [Features](#features)
   - [Technologies](#Technologies)
2. [Start](#start)
   - [Prerequisite](#Prerequisite)
   - [Installation](#Installation)
   - [Configuration](#Configuration)
   - [Run](#Run)
   - [Run via Container](#RunViaContainer)
3. [Contact](#Contact)

## About

This Coffee Shop website was built with the aim of making it easier for users to manage a coffee shop business. On the Backend, this website uses Golang with the Gin Gonic framework, and the database uses PostgreSQL.

### Features

- CRUD User, Product, Favorite, Order
- Authentication With JWT
- Hash Password
- Cloudinary
- Seeding

### Technologies

- Gin Gonic
- Golang
- PostgreSQL

## Start

### Prerequisite

To get started, you need to have Golang installed on your system. If it's not installed yet, download and install it from the official Golang website.

### Installation

1. Clone the repository

```sh
$ git clone https://github.com/mfauzidr/Coffee-Shop-Go-Backend.git
```

2. Download the dependencies:

```sh
$ go mod tidy
```

### Configuration

The project uses a .env file for environment variables like database connection details, server port, etc.
you can create a .env file according to the .env.example in the root directory

### Run

Run the following command to start the server:

```sh
$ go run ./cmd/main.go
```


## Contact

mfdr.fauzi97@gmail.com || Mochammad Fauzi Dwi R
Github : https://github.com/mfauzidr/Coffee-Shop-Go-Backend
