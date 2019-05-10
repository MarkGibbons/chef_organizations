CREATE DATABASE IF NOT EXISTS organizations;

USE organizations;

CREATE TABLE IF NOT EXISTS organizations (
  name TEXT,
  full_name  TEXT
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS org_groups (
  name TEXT,
  group_name  TEXT,
  organization_name  TEXT,
  user_name  TEXT,
  actor  BOOLEAN,
  user  BOOLEAN
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS members (
  user_name TEXT,
  email  TEXT,
  display_name  TEXT
) ENGINE=INNODB;
