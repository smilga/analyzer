/*
delimiter //
CREATE TRIGGER archive_matches
AFTER DELETE
  ON matches FOR EACH ROW
BEGIN
  INSERT INTO matches_archive (pattern_id, website_id, report_id, value, created_at, deleted_at)
  VALUES (old.pattern_id, old.website_id, old.report_id, old.value, old.created_at, NOW());
END ; //
delimiter ;
*/
