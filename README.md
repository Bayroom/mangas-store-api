# Manga Store API

This is a simple CRUD API for managing manga books in a store. The API is built using Go (Gin Gonic) and can be containerized with Docker.

## API Endpoints

- `GET /`: Welcome message
- `GET /mangas`: Get all mangas
- `GET /mangas/:id`: Get a specific manga by ID
- `POST /mangas`: Add a new manga
- `PUT /mangas/:id`: Update a manga by ID
- `DELETE /mangas/:id`: Delete a manga by ID

## How to Run Locally

1. Install Go and Docker
2. Clone the repository: `git clone https://github.com/Bayroom/mangas-store-api.git`
3. Navigate to the project directory: `cd mangas-store-api`
4. Build the Docker image: `docker build -t mangas-store-api .`
5. Run the Docker container: `docker run -p 8080:8080 mangas-store-api`
6. Access the API at `http://localhost:8080`

## Contributing

Feel free to contribute by opening issues or pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.