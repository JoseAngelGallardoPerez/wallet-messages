<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class InitTables extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
    }

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        // skip the migration if there are another migrations
        // It means this migration was already applied
        $migrations = DB::select('SELECT * FROM migrations LIMIT 1');
        if (!empty($migrations)) {
            return;
        }
        $oldMigrationTable = DB::select("SHOW TABLES LIKE 'schema_migrations'");
        if (!empty($oldMigrationTable)) {
            return;
        }

        DB::beginTransaction();

        try {
            app("db")->getPdo()->exec($this->getSql());
        } catch (\Throwable $e) {
            DB::rollBack();
            throw $e;
        }

        DB::commit();
    }

    private function getSql()
    {
        return <<<SQL
            CREATE TABLE `messages` (
              `id` int(11) UNSIGNED NOT NULL,
              `message` text NOT NULL,
              `subject` varchar(255) DEFAULT NULL,
              `recipient_id` varchar(36) DEFAULT NULL,
              `sender_id` varchar(36) DEFAULT NULL,
              `edited` tinyint(1) DEFAULT NULL,
              `created_at` timestamp NULL DEFAULT NULL,
              `updated_at` timestamp NULL DEFAULT NULL,
              `deleted_for_sender` tinyint(4) DEFAULT NULL,
              `deleted_for_recipient` tinyint(4) DEFAULT NULL,
              `parent_id` int(11) UNSIGNED DEFAULT NULL,
              `is_sender_read` tinyint(4) DEFAULT NULL,
              `is_recipient_read` tinyint(4) DEFAULT NULL,
              `delete_after_read` tinyint(1) DEFAULT '0',
              `is_recipient_incoming` tinyint(1) DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            CREATE TABLE `schema_migrations` (
              `version` bigint(20) NOT NULL,
              `dirty` tinyint(1) NOT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `schema_migrations` (`version`, `dirty`) VALUES
            (20180926221121, 0);


            ALTER TABLE `messages`
              ADD PRIMARY KEY (`id`),
              ADD KEY `parent_index` (`parent_id`);

            ALTER TABLE `schema_migrations`
              ADD PRIMARY KEY (`version`);


            ALTER TABLE `messages`
              MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;


            ALTER TABLE `messages`
              ADD CONSTRAINT `FK_messages_parent` FOREIGN KEY (`parent_id`) REFERENCES `messages` (`id`);
SQL;
    }
}
