CREATE TABLE IF NOT EXISTS matches (
    id INT NOT NULL AUTO_INCREMENT,
    pattern_id INT NOT NULL,
    website_id INT NOT NULL,
    report_id INT NOT NULL,
    value VARCHAR(500) NOT NULL,
    created_at DATETIME NULL DEFAULT NULL,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (pattern_id) REFERENCES patterns(id),
    FOREIGN KEY (website_id) REFERENCES websites(id),
    FOREIGN KEY (report_id) REFERENCES reports(id)
) DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
