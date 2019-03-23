CREATE TABLE IF NOT EXISTS reports (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    website_id INT NOT NULL,
    loaded_in VARCHAR(10) NULL DEFAULT NULL,
    resource_check_in VARCHAR(10) NULL DEFAULT NULL,
    html_check_in VARCHAR(10) NULL DEFAULT NULL,
    total_in VARCHAR(10) NULL DEFAULT NULL,
    created_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (website_id) REFERENCES websites(id)
) DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
