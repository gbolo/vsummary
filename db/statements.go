package db

const (

	// PREPARED INSERT STATEMENTS --------------------------------------------------------------------------------------

	insertPoller = `
INSERT INTO poller (
	vcenter_host,
	vcenter_name,
	enabled,
	user_name,
	password,
	interval_min
	)
VALUES (
	:vcenter_host,
	:vcenter_name,
	:enabled,
	:user_name,
	:password,
	:interval_min
	)
ON DUPLICATE KEY UPDATE
	vcenter_name=VALUES(vcenter_name),
	enabled=VALUES(enabled),
	user_name=VALUES(user_name),
	password=VALUES(password),
	interval_min=VALUES(interval_min);`
)
