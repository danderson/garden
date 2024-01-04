-- First delete tables with foreign keys
drop table auth_group_permissions;
drop table auth_user_groups;
drop table auth_user_user_permissions;
drop table django_admin_log;
drop table auth_permission;

-- Then, the rest.
drop table django_migrations;
drop table django_content_type;
drop table auth_group;
drop table auth_user;
drop table django_session;
