# EZ Tools

# UI - ui.go

Debugging and Verbose to control interaction and verbose levels.

# Xml - cfg.go

It reads from an xml into a structure.

# App self-update - appup.go

It gets server info from database and updates in background. The next time the app is run, it will be in the new version.

# LDAP - auth.go

It gets server info, including DN structure, from database and tries to bind to check for authentication.

# DB - conn.go, tbl.go, pair.go

It provides operations to connect, insert, search from databases. pair.go is for id(int)-str(string) tables, while tbl.go is for more various types of tables.

# Restful APIs - rest.go

Use this to send and receive through Restful APIs. Based on the needs for Jira and Gerrit.

# Contacts - teams/, contacts

It provides operations to manage contacts and teams in databases.

# tables

Miscellaneous
CREATE TABLE chore (id CHAR(50) NOT NULL UNIQUE, str TINYTEXT, PRIMARY KEY(id)); 

For Google tests
CREATE TABLE google (id TINYINT AUTO_INCREMENT, tool TINYINT, android TINYINT, ver TINYINT, req DATE, exp DATE, PRIMARY KEY(id)); 
CREATE TABLE prodgle (google TINYINT, product TINYINT); 
CREATE TABLE tool (id TINYINT AUTO_INCREMENT, str TINYTEXT, PRIMARY KEY(id)); 
CREATE TABLE product (id TINYINT, str TINYTEXT); 
CREATE TABLE bit (id TINYINT AUTO_INCREMENT, str TINYTEXT, PRIMARY KEY(id)); 
CREATE TABLE phase (id TINYINT AUTO_INCREMENT, str TINYTEXT, PRIMARY KEY(id)); 
CREATE TABLE android (id TINYINT AUTO_INCREMENT, str TINYTEXT, PRIMARY KEY(id)); 
CREATE TABLE prodfo (product TINYINT, bit TINYINT, android TINYINT, phase TINYINT); 
CREATE TABLE ver (id TINYINT AUTO_INCREMENT, str TINYTEXT, PRIMARY KEY(id)); 

For SOC
CREATE TABLE carrier (id TINYINT, str TINYTEXT);
CREATE TABLE soc (carrier TINYTEXT, product TINYTEXT, android TINYTEXT, gms TINYTEXT, section SMALLINT, subsection SMALLINT, identifier TINYTEXT, description VARCHAR(200), team TINYINT, compliance TINYINT, comments VARCHAR(100)); 
CREATE TABLE sec (id SMALLINT, str TINYTEXT); 
CREATE TABLE subsec (id SMALLINT, str TINYTEXT);
CREATE TABLE compl (id TINYINT, str TINYTEXT);

For contacts
CREATE TABLE team (id TINYINT AUTO_INCREMENT, str TINYTEXT, leader TINYINT, PRIMARY KEY(id)); 
CREATE TABLE contacts (id TINYINT AUTO_INCREMENT, number TINYTEXT, name TINYTEXT NOT NULL, team TINYINT, ext TINYTEXT, phone TINYTEXT, mail TINYTEXT NOT NULL, ldap TINYTEXT, uid TINYTEXT, PRIMARY KEY(id)) CHARSET=utf8;

For weekly reports
CREATE TABLE weeklyTaskBars (id TINYINT NOT NULL UNIQUE, str TINYTEXT NOT NULL, PRIMARY KEY(id)); 
CREATE TABLE weeklyTaskTitles (id TINYINT AUTO_INCREMENT, str TINYTEXT NOT NULL, PRIMARY KEY(id)); 
CREATE TABLE weeklyTaskCurr (id TINYINT NOT NULL UNIQUE, str TINYINT NOT NULL UNIQUE); 
CREATE TABLE weeklyTaskNext (id TINYINT NOT NULL UNIQUE, str TINYINT NOT NULL UNIQUE); 
CREATE TABLE weeklyTaskDesc (id TINYINT AUTO_INCREMENT, str VARCHAR(500) NOT NULL UNIQUE, PRIMARY KEY(id));
CREATE TABLE weeklyTaskWork (id TINYINT AUTO_INCREMENT, str TINYINT NOT NULL, contact TINYINT NOT NULL, section TINYINT NOT NULL, week TINYINT NOT NULL, PRIMARY KEY(id)); 

| table            | field       | comment                                                      | example            |
| :--------------- | :---------- | :----------------------------------------------------------- | :----------------- |
| chore            | id          | a string for name                                            | googleVer          |
|                  | str         | a string for value                                           | 1.0                |
| sec              | id          | one id may relate to multiple strings                        | 5                  |
|                  | str         |                                                              | Browser            |
| subsec           | id          | one id may relate to multiple strings                        | 1                  |
|                  | str         |                                                              | WAP                |
| team             | id          |                                                              | 2                  |
|                  | str         |                                                              | BSP                |
|                  | leader      | id in table contacts                                         | 16                 |
| compl            | id          | one id may relate to multiple strings                        | 3                  |
|                  | str         |                                                              | Complies           |
| carrier          | id          | one id may relate to multiple strings                        | 4                  |
|                  | str         |                                                              | CMCC               |
| tool             | id          |                                                              | 10                 |
|                  | str         |                                                              | CTS                |
| ver              | id          |                                                              | 13                 |
|                  | str         |                                                              | 201805             |
| product          | id          | one id may relate to multiple strings                        | 11                 |
|                  | str         |                                                              | Pixel              |
| bit              | id          |                                                              | 15                 |
|                  | str         |                                                              | 32/64/go           |
| prodfo           | product     |                                                              | 11                 |
|                  | bit         |                                                              | 15                 |
|                  | android     |                                                              | 12                 |
|                  | phase       |                                                              | 16                 |
| phase            | id          |                                                              | 16                 |
|                  | str         |                                                              | MP                 |
| prodgle          | google      |                                                              | 14                 |
|                  | product     |                                                              | 11                 |
| android          | id          | if all wanted, apply to all current                          | 12                 |
|                  | str         |                                                              | 8.1                |
| google           | id          |                                                              | 14                 |
|                  | tool        | id in table tool                                             | 10                 |
|                  | android     | id in table android                                          | 12                 |
|                  | ver         | id in table ver                                              | 13                 |
|                  | req         | date of requirement                                          | 20080401           |
|                  | exp         | date of expiry                                               | 20080501           |
| soc              | carrier     | string                                                       | CMCC               |
|                  | product     | string                                                       | Pixel              |
|                  | android     | string                                                       | 8.1Go              |
|                  | gms         | string                                                       | 201805             |
|                  | section     | id in table sec                                              | 5                  |
|                  | subsection  | id in table subsec                                           | 1                  |
|                  | identifier  | string                                                       | MSR19820           |
|                  | description | string                                                       | detail..           |
|                  | team        | id in table team                                             | 2                  |
|                  | compliance  | id in table compl                                            | 3                  |
|                  | comments    | string                                                       | not supp           |
| contacts         | id          |                                                              | 16                 |
|                  | number      | string                                                       | 007                |
|                  | name        | string                                                       | Allen              |
|                  | team        | id in table team                                             | 2                  |
|                  | ext         | string                                                       | 12345              |
|                  | phone       | string                                                       | 1234567890         |
|                  | mail        | string                                                       | allen              |
|                  | ldap        | string of user name for LDAP commands                        | Allen              |
|                  | uid         | string of user name for code servers                         | allen              |
| weeklyTaskBars   | id          | 0, 1 for current & next week                                 | 23                 |
|                  | str         | string                                                       | Tasks Accomplished |
| weeklyTaskTitles | id          |                                                              | 17                 |
|                  | str         | string                                                       | Project 1          |
| weeklyTaskCurr   | id          | the tens place is the page number of ppt. the ones place is the order within a page. | 18                 |
|                  | str         | id in table weeklyTaskTitles                                 | 17                 |
| weeklyTaskNext   | id          | the tens place is the page number of ppt. the ones place is the order within a page. | 19                 |
|                  | str         | id in table weeklyTaskTitles                                 | 17                 |
| weeklyTaskDesc   | id          |                                                              | 20                 |
|                  | str         | unique work description                                      | bug fix            |
| weeklyTaskWork   | id          |                                                              | 21                 |
|                  | str         | id in weeklyTaskDesc                                         | 20                 |
|                  | contact     | id in table contacts                                         | 16                 |
|                  | section     | id in table weeklyTaskCurr/Next                              | 18                 |
|                  | week        | week number as int                                           | 22                 |



## table details

chore's ids

- GoogleUrl: Url of google tool app
- GoogleUrlDev: Dev Url of google tool app
- GoogleDir: dir under Url
- GoogleApp: app name of google tool app
- GoogleRes: resource request name under Url
- GooglexxxFolder: root path on the remote or local server. xxx=Svr/Lcl
- GooglexxxPath_XXX: tool paths on the servers. XXX=CTS/ETS/GTS/STS/VTS/PAB/GMS
- GooglexxxPath_Android: prefix to Android versions of the paths on the servers
- GooglexxxSubAndroid: Among CTS, ETS, GTS, STS, VTS, PAB & GMS, which have sub-directories for Android versions, appended by a space for each to separate.
- GoogleSvrSTS_pre: prefix of files
- GoogleSvrVTS_pre: prefix of files all bits need
- GoogleSvrPAB_scr: prefix of script files all bits need
- GoogleSvrPAB_img: prefix of image files different bits need
- LogTran: log file name for file transport. comma separated csv format with no headers and columns are source file name and destination file name
- LogSize: log file size. usage reserved
- WeeklyPptUrl: Url of weekly report app
- WeeklyPptDir: dir under Url
- WeeklyPptApp: app name of weekly report tool app
- WeeklyPptUrlDev: dev url of weekly report tool app
- WeeklyPptStaticPath: dir of PPT on static server
- WeeklyPptStaticSuf: suffix of PPT on static server. defaults to .pptx
- AuthIP: IP of authentication server
- AuthRootD: former part of root bind dn before user name
- AuthRootN: latter part of root bind dn after user name
