alter table screens_items add	rowspan		int(4)	DEFAULT '0' NOT NULL;

drop table escalation_rules;
drop table escalations;

--
-- Table structure for table 'escalations'
--

CREATE TABLE escalations (
  escalationid		int(4)		NOT NULL auto_increment,
  name			varchar(64)	DEFAULT '0' NOT NULL,
  dflt			int(2)		DEFAULT '0' NOT NULL,
  PRIMARY KEY (escalationid),
  UNIQUE (name)
) type=InnoDB;

--
-- Table structure for table 'escalation_rules'
--

CREATE TABLE escalation_rules (
  escalationruleid		int(4)		NOT NULL auto_increment,
  escalationid		int(4)		DEFAULT '0' NOT NULL,
  level			int(4)		DEFAULT '0' NOT NULL,
  period		varchar(100)	DEFAULT '1-7,00:00-23:59' NOT NULL,
  delay			int(4)		DEFAULT '0' NOT NULL,
  actiontype		int(4)		DEFAULT '0' NOT NULL,
  PRIMARY KEY (escalationruleid),
  KEY (escalationid)
) type=InnoDB;
