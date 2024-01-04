-- Django tables
CREATE TABLE IF NOT EXISTS django_migrations (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       app varchar(255) NOT NULL,
       name varchar(255) NOT NULL,
       applied datetime NOT NULL);
CREATE TABLE IF NOT EXISTS auth_group_permissions (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       group_id integer NOT NULL REFERENCES auth_group (id) DEFERRABLE INITIALLY DEFERRED,
       permission_id integer NOT NULL REFERENCES auth_permission (id) DEFERRABLE INITIALLY DEFERRED);
CREATE TABLE IF NOT EXISTS auth_user_groups (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       user_id integer NOT NULL REFERENCES auth_user (id) DEFERRABLE INITIALLY DEFERRED,
       group_id integer NOT NULL REFERENCES auth_group (id) DEFERRABLE INITIALLY DEFERRED);
CREATE TABLE IF NOT EXISTS auth_user_user_permissions (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       user_id integer NOT NULL REFERENCES auth_user (id) DEFERRABLE INITIALLY DEFERRED,
       permission_id integer NOT NULL REFERENCES auth_permission (id) DEFERRABLE INITIALLY DEFERRED);
CREATE TABLE IF NOT EXISTS django_admin_log (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       object_id text,
       object_repr varchar(200) NOT NULL,
       action_flag smallint NOT NULL CHECK (action_flag >= 0),
       change_message text NOT NULL,
       content_type_id integer REFERENCES django_content_type (id) DEFERRABLE INITIALLY DEFERRED,
       user_id integer NOT NULL REFERENCES auth_user (id) DEFERRABLE INITIALLY DEFERRED,
       action_time datetime NOT NULL);
CREATE TABLE IF NOT EXISTS django_content_type (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       app_label varchar(100) NOT NULL,
       model varchar(100) NOT NULL);
CREATE TABLE IF NOT EXISTS auth_permission (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       content_type_id integer NOT NULL REFERENCES django_content_type (id) DEFERRABLE INITIALLY DEFERRED,
       codename varchar(100) NOT NULL, name varchar(255) NOT NULL);
CREATE TABLE IF NOT EXISTS auth_group (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       name varchar(150) NOT NULL UNIQUE);
CREATE TABLE IF NOT EXISTS auth_user (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       password varchar(128) NOT NULL,
       last_login datetime,
       is_superuser bool NOT NULL,
       username varchar(150) NOT NULL UNIQUE,
       last_name varchar(150) NOT NULL,
       email varchar(254) NOT NULL,
       is_staff bool NOT NULL,
       is_active bool NOT NULL,
       date_joined datetime NOT NULL,
       first_name varchar(150) NOT NULL);
CREATE TABLE IF NOT EXISTS django_session (
       session_key varchar(40) NOT NULL PRIMARY KEY,
       session_data text NOT NULL,
       expire_date datetime NOT NULL);

CREATE UNIQUE INDEX auth_group_permissions_group_id_permission_id_0cd325b0_uniq ON auth_group_permissions (group_id, permission_id);
CREATE INDEX auth_group_permissions_group_id_b120cbf9 ON auth_group_permissions (group_id);
CREATE INDEX auth_group_permissions_permission_id_84c5c92e ON auth_group_permissions (permission_id);
CREATE UNIQUE INDEX auth_user_groups_user_id_group_id_94350c0c_uniq ON auth_user_groups (user_id, group_id);
CREATE INDEX auth_user_groups_user_id_6a12ed8b ON auth_user_groups (user_id);
CREATE INDEX auth_user_groups_group_id_97559544 ON auth_user_groups (group_id);
CREATE UNIQUE INDEX auth_user_user_permissions_user_id_permission_id_14a6b632_uniq ON auth_user_user_permissions (user_id, permission_id);
CREATE INDEX auth_user_user_permissions_user_id_a95ead1b ON auth_user_user_permissions (user_id);
CREATE INDEX auth_user_user_permissions_permission_id_1fbb5f2c ON auth_user_user_permissions (permission_id);
CREATE INDEX django_admin_log_content_type_id_c4bce8eb ON django_admin_log (content_type_id);
CREATE INDEX django_admin_log_user_id_c564eba6 ON django_admin_log (user_id);
CREATE UNIQUE INDEX django_content_type_app_label_model_76bd3d3b_uniq ON django_content_type (app_label, model);
CREATE UNIQUE INDEX auth_permission_content_type_id_codename_01ab375a_uniq ON auth_permission (content_type_id, codename);
CREATE INDEX auth_permission_content_type_id_2f476e4b ON auth_permission (content_type_id);
CREATE INDEX django_session_expire_date_a5c62663 ON django_session (expire_date);

-- Herbarium tables
CREATE TABLE IF NOT EXISTS herbarium_family (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       name varchar(200) NOT NULL);
CREATE TABLE IF NOT EXISTS herbarium_plantname (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       name varchar(200) NOT NULL,
       plant_id bigint NOT NULL REFERENCES herbarium_plant (id) DEFERRABLE INITIALLY DEFERRED);
CREATE TABLE IF NOT EXISTS herbarium_plant (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       edible bool,
       needs_trellis bool,
       needs_bird_netting bool,
       is_keto bool,
       native bool,
       invasive bool,
       is_cover bool,
       grow_from_seed bool,
       type varchar(1) NOT NULL,
       lifespan varchar(1) NOT NULL,
       bad_for_cats bool,
       deer_resistant bool,
       family_id bigint NOT NULL REFERENCES herbarium_family (id) DEFERRABLE INITIALLY DEFERRED,
       name varchar(200) NOT NULL);
CREATE TABLE IF NOT EXISTS herbarium_tag (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       name varchar(200) NOT NULL UNIQUE);
CREATE TABLE IF NOT EXISTS herbarium_plant_tags (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       plant_id bigint NOT NULL REFERENCES herbarium_plant (id) DEFERRABLE INITIALLY DEFERRED,
       tag_id bigint NOT NULL REFERENCES herbarium_tag (id) DEFERRABLE INITIALLY DEFERRED);
CREATE TABLE IF NOT EXISTS herbarium_plantingwindow (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       type varchar(1) NOT NULL,
       end smallint NOT NULL,
       start smallint NOT NULL,
       plant_id bigint NOT NULL REFERENCES herbarium_plant (id) DEFERRABLE INITIALLY DEFERRED);
CREATE TABLE IF NOT EXISTS boxinventory_boxcontent (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       name varchar(200) NOT NULL,
       box_id bigint NOT NULL REFERENCES boxinventory_box (id) DEFERRABLE INITIALLY DEFERRED,
       planted date NOT NULL, removed date);
CREATE TABLE IF NOT EXISTS boxinventory_box (
       id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
       name varchar(200) NOT NULL,
       want_qr bool NOT NULL,
       qr_applied bool NOT NULL);

CREATE INDEX herbarium_plantname_plant_id_558f61f2 ON herbarium_plantname (plant_id);
CREATE INDEX herbarium_plant_family_id_dfa6b981 ON herbarium_plant (family_id);
CREATE UNIQUE INDEX herbarium_plant_tags_plant_id_tag_id_cdbc19db_uniq ON herbarium_plant_tags (plant_id, tag_id);
CREATE INDEX herbarium_plant_tags_plant_id_8a1d0e42 ON herbarium_plant_tags (plant_id);
CREATE INDEX herbarium_plant_tags_tag_id_ae4f410c ON herbarium_plant_tags (tag_id);
CREATE INDEX herbarium_plantingwindow_plant_id_00a85eb3 ON herbarium_plantingwindow (plant_id);
CREATE INDEX boxinventory_boxcontent_box_id_e224e154 ON boxinventory_boxcontent (box_id);
