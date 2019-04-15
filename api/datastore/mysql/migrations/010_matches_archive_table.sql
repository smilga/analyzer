CREATE TABLE IF NOT EXISTS matches_archive (
    id INT NOT NULL AUTO_INCREMENT,
    pattern_id INT NOT NULL,
    website_id INT NOT NULL,
    report_id INT NOT NULL,
    value TEXT NOT NULL,
    created_at DATETIME NULL DEFAULT NULL,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
) DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
