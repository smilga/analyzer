CREATE TABLE IF NOT EXISTS filter_tags (
    id INT NOT NULL AUTO_INCREMENT,
    filter_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (filter_id) REFERENCES filters(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
) DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
