CREATE TABLE IF NOT EXISTS pattern_tags (
    id INT NOT NULL AUTO_INCREMENT,
    pattern_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (pattern_id) REFERENCES patterns(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
) DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
