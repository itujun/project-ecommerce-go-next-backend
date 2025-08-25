-- refresh_tokens: menyimpan RT yang berlaku (opsional multi-device)
CREATE TABLE IF NOT EXISTS refresh_tokens (
  id           CHAR(36)      NOT NULL PRIMARY KEY,     -- jti (UUID)
  user_id      CHAR(36)      NOT NULL,                 -- relasi ke users.id
  token_hash   CHAR(64)      NOT NULL,                 -- sha256 hash dari RT
  issued_at    DATETIME      NOT NULL,
  expires_at   DATETIME      NOT NULL,
  revoked      TINYINT(1)    NOT NULL DEFAULT 0,
  created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_rt_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_rt_user ON refresh_tokens(user_id);
CREATE INDEX idx_rt_expires ON refresh_tokens(expires_at);
