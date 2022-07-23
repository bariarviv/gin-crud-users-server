CREATE TABLE IF NOT EXISTS users(
    email            VARCHAR(200) PRIMARY KEY,
    username         VARCHAR(50) NOT NULL,
    password         VARCHAR(50) NOT NULL,
    sys_created_date TIMESTAMP with time zone NOT NULL DEFAULT now()
);
