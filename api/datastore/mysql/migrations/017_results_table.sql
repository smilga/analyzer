CREATE TABLE IF NOT EXISTS service_features (
    id INT NOT NULL AUTO_INCREMENT,
    service_id INT NOT NULL,
    feature_id INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (service_id) REFERENCES services(id),
    FOREIGN KEY (feature_id) REFERENCES features(id)
) DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
