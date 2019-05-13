CREATE DATABASE IF NOT EXISTS organizations;

USE organizations;

CREATE TABLE IF NOT EXISTS organizations (
  name TEXT,
  full_name  TEXT
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS org_groups (
  name TEXT,
  organization_name  TEXT,
  user_name  TEXT,
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS members (
  user_name TEXT,
  email  TEXT,
  display_name  TEXT
) ENGINE=INNODB;
