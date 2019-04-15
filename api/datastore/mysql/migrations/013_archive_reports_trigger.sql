delimiter //
CREATE TRIGGER archive_reports
AFTER DELETE
  ON reports FOR EACH ROW
BEGIN
  INSERT INTO reports_archive (user_id, website_id, started_in, loaded_in, resource_check_in, html_check_in, total_in, created_at, deleted_at)
  VALUES (old.user_id, old.website_id, old.started_in, old.loaded_in, old.resource_check_in, old.html_check_in, old.total_in, old.created_at, NOW());
END ; //
delimiter ;
