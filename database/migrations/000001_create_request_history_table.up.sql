CREATE TABLE request_history (
  request_id SMALLINT UNSIGNED AUTO_INCREMENT,
  endpoint VARCHAR(1000) NOT NULL,
  headers VARCHAR(2000),
  method ENUM("GET", "POST", "TRACE", "HEAD", "DELETE", "PATCH", "PUT"),
  body VARCHAR(2000),
  request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT pk_history PRIMARY KEY (request_id)
);
