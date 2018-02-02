package db

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"fmt"

	"github.com/gbolo/vsummary/common"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var (

	insertVCenter = `
INSERT INTO vcenter (id, fqdn, short_name, user_name, password)
	VALUES (:id, :fqdn, :short_name, :user_name, :password);`

	SchemaVCenter = `
CREATE TABLE IF NOT EXISTS vcenter (
    id text,
    fqdn text,
    short_name text,
	user_name text,
	password text);`

	insertVm = `
INSERT INTO vm (
	id,
	name,
	moref,
	vmx_path,
	vcpu,
	memory_mb,
	template,
	config_guest_os,
	config_version,
	smbios_uuid,
	instance_uuid,
	config_change_version,
	guest_tools_version,
	guest_tools_running,
	guest_hostname,
	guest_ip,
	guest_os,
	stat_cpu_usage,
	stat_host_memory_usage,
	stat_guest_memory_usage,
	stat_uptime_sec,
	power_state,
	folder_id,
	vapp_id,
	resourcepool_id,
	esxi_id,
	vcenter_id
	)
VALUES (
	:id,
	:name,
	:moref,
	:vmx_path,
	:vcpu,
	:memory_mb,
	:template,
	:config_guest_os,
	:config_version,
	:smbios_uuid,
	:instance_uuid,
	:config_change_version,
	:guest_tools_version,
	:guest_tools_running,
	:guest_hostname,
	:guest_ip,
	:guest_os,
	:stat_cpu_usage,
	:stat_host_memory_usage,
	:stat_guest_memory_usage,
	:stat_uptime_sec,
	:power_state,
	:folder_id,
	:vapp_id,
	:resourcepool_id,
	:esxi_id,
	:vcenter_id
	)
ON DUPLICATE KEY UPDATE
	id=VALUES(id),
	name=VALUES(name),
	moref=VALUES(moref),
	vmx_path=VALUES(vmx_path),
	vcpu=VALUES(vcpu),
	memory_mb=VALUES(memory_mb),
	template=VALUES(template),
	config_guest_os=VALUES(config_guest_os),
	config_version=VALUES(config_version),
	smbios_uuid=VALUES(smbios_uuid),
	instance_uuid=VALUES(instance_uuid),
	config_change_version=VALUES(config_change_version),
	guest_tools_version=VALUES(guest_tools_version),
	guest_tools_running=VALUES(guest_tools_running),
	guest_hostname=VALUES(guest_hostname),
	guest_ip=VALUES(guest_ip),
	guest_os=VALUES(guest_os),
	stat_cpu_usage=VALUES(stat_cpu_usage),
	stat_host_memory_usage=VALUES(stat_host_memory_usage),
	stat_guest_memory_usage=VALUES(stat_guest_memory_usage),
	stat_uptime_sec=VALUES(stat_uptime_sec),
	power_state=VALUES(power_state),
	folder_id=VALUES(folder_id),
	vapp_id=VALUES(vapp_id),
	resourcepool_id=VALUES(resourcepool_id),
	esxi_id=VALUES(esxi_id),
	vcenter_id=VALUES(vcenter_id),
	present=1
;
`

)

var log = logging.MustGetLogger("vsummary")

// UserRecord defines the properties of a user
type VCenterRecord struct {
	Id             string `db:"id"`
	FQDN           string `db:"fqdn"`
	ShortName      string `db:"short_name"`
	UserName       string `db:"user_name"`
	Password       string `db:"password"`
	Other		   string
}

// Backend interface.
type Backend struct {
	db *sqlx.DB
}

// new Backend
func InitBackend() (b *Backend, err error) {

	driver := viper.GetString("backend.db_driver")
	dsn := viper.GetString("backend.db_dsn")

	db, err := sqlx.Connect(driver, dsn)

	if err != nil {
		log.Error("failed to connect to database")
		err = fmt.Errorf("database error: %s", err)
	} else {
		b = &Backend{ db: db }
		log.Info("connection to database successful")
	}

	return
}

func NewBackend() *Backend {
	return &Backend{}
}

// Check DB connection
func (b *Backend) checkDB() error {
	if b.db == nil {
		return errors.New("Database connection is not set up")
	}
	return nil
}

// SetDB changes the underlying sql.DB object Accessor is manipulating.
func (b *Backend) SetDB(db *sqlx.DB) {
	b.db = db
}


func (b *Backend) ApplySchemas() (err error) {

	// check if db connection is available
	err = b.checkDB()
	if err != nil {
		return
	}

	// apply all schemas
	for _, schema := range generateSqlSchemas() {

		log.Debugf("Applying schema: %s", schema.Name)
		_, err = b.db.MustExec(schema.SqlQuery).RowsAffected()
		if err == nil {
			log.Debugf("db schema is present: %s", schema.Name)
		} else {
			break
		}

	}

	return
}

// vCenterRecord inserts user into database
func (b *Backend) InsertVCenter(vcenter *VCenterRecord) error {

	err := b.checkDB()
	if err != nil {
		return err
	}


	// Store the user record in the DB
	res, err := b.db.NamedExec(insertVCenter, &vcenter)

	if err != nil {
		//log.Errorf("Error adding identity %s to the database: %s", user.Name, err)
		return err
	}

	numRowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if numRowsAffected == 0 {
		return fmt.Errorf("Failed to add identity %s to the database", vcenter.FQDN)
	}

	if numRowsAffected != 1 {
		return fmt.Errorf("Expected to add one record to the database, but %d records were added", numRowsAffected)
	}

	return nil

}

// InsertVMs inserts a vm into database
func (b *Backend) InsertVMs(vms []common.Vm) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this
	tx.MustExec("UPDATE vm SET present = 0 WHERE present = 1 AND vcenter_id=?", vms[0].VcenterId)
	var rowsAffected int64 = 0

	for _, vm := range vms {
		// fill in missing data
		// folder may not exist
		if vm.FolderMoref != "" && vm.FolderMoref != "vapp" {
			vm.FolderId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.FolderMoref))
		}

		// vapps may not exist
		if vm.VappMoref != "none" {
			vm.VappId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.VappMoref))
		}
		// resourcepool may not exist
		if vm.ResourcePoolMoref != "" {
			vm.ResourcePoolId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.ResourcePoolId))
		}

		// Fill in some required Ids
		vm.Id = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.Moref))
		vm.EsxiId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.EsxiMoref))

		// Store the user record in the DB
		res, err := tx.NamedExec(insertVm, &vm)

		if err != nil {
			break
		}

		numRowsAffected, err := res.RowsAffected()
		if err != nil {
			break
		}

		rowsAffected = rowsAffected + numRowsAffected

	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit transaction to database: %s", err)
	}

	log.Debugf("total combined affected rows: %d", rowsAffected)

	return

}