-- Create "meeting_types" table
CREATE TABLE `meeting_types` (`id` text NOT NULL, `name` text NOT NULL, `name_key` text NOT NULL, `duration_minutes` integer NOT NULL, `description` text NULL, PRIMARY KEY (`id`));
-- Create index "meeting_types_name_key_key" to table: "meeting_types"
CREATE UNIQUE INDEX `meeting_types_name_key_key` ON `meeting_types` (`name_key`);
