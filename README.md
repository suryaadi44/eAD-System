# Electronic Administration Documentation System

Simple RESTful API for Administration Documentation System build on Go and Echo Framework.

This project is a part of miniprojects for Alterra Become a Master of Golang Programming Course ft. Kampus Merdeka Program

By : Surya Adi (kmsurya.adi44@gmail.com)

## Running Pre-requisites

Before running the application, you need to install the following:

-   Go 1.19
-   MySQL 8.0
-   [wkhtmltopdf 0.12.5 (with patched qt)](https://wkhtmltopdf.org/downloads.html)

Or you can use docker to run the application by building image from Dockerfile or using this image `suryawarior44/ead-system` from docker hub.

## Running the Application

This app need to run with environment variables, you can set the environment variables in `.env` file or set it in your OS environment variables.

Environment variables needed:

| Name        | Description                                                                      |
| ----------- | -------------------------------------------------------------------------------- |
| DB_HOST     | Database host                                                                    |
| DB_PORT     | Database port                                                                    |
| DB_USER     | Database username                                                                |
| DB_PASSWORD | Database password                                                                |
| DB_NAME     | Database name                                                                    |
| PORT        | Port for the application to run on                                               |
| JWT_SECRET  | Secret key for JWT                                                               |
| QR_PATH     | URL Path for document checking endpoint , eg: `http://localhost:8080/documents/` |

