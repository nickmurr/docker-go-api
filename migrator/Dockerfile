FROM migrate/migrate

COPY /migrations /migrations

CMD ["-path", "/migrations", "-database", "postgresql://postgres:docker@postgres:5432/postgres?sslmode=disable", "up"]
CMD ["-path", "/migrations", "-database", "postgresql://postgres:docker@postgres:5432/restapi_test?sslmode=disable", "up"]
