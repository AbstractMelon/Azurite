-- Add new fields to game_requests table
ALTER TABLE game_requests ADD COLUMN reason TEXT;
ALTER TABLE game_requests ADD COLUMN existing_community TEXT;
ALTER TABLE game_requests ADD COLUMN mod_loader TEXT;
ALTER TABLE game_requests ADD COLUMN contact TEXT;