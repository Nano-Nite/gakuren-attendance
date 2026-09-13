package helper

import (
	"log"
	"os"

	"github.com/google/uuid"
)

var API_VERSION = os.Getenv("API_VERSION")

// COMMON CONST
const REJECT = "REJECT"
const SAVE = "SAVE"
const WORKFLOW_NOTFOUND_BEHAVIOUR = "workflow.notfound.behaviour"
const APPROVAL_DOCUMENT_CODE = "APR"

// ENTITY TYPE
const CLASS_ENTITY_TYPE = "CLASS"
const USER_ENTITY_TYPE = "USER"
const STUDENT_ENTITY_TYPE = "STUDENT"
const TNS_ENTITY_TYPE = "TNS"
const LOCATION_ENTITY_TYPE = "LOCATION"

// MODULE CODE
const CLASS_MODULE_CODE = "KLS"
const STUDENT_MODULE_CODE = "MRD"
const TNS_MODULE_CODE = "TNS"
const LOCATION_MODULE_CODE = "LOC"

// DEFAULT SETTING
const DEFAULT_ROW_PER_PAGES = 10
const DEFAULT_PAGES = 1

// STATUS NAME
const STATUS_ACTIVE = "Active"
const STATUS_INACTIVE = "Inactive"
const STATUS_PENDING = "Pending"
const STATUS_DELETED = "Deleted"
const STATUS_NEW_USER = "New User"

// ACTION CODE
const ACTION_CODE_CREATE = "CREATE"
const ACTION_CODE_UPDATE = "UPDATE"
const ACTION_CODE_DELETE = "DELETE"
const ACTION_CODE_ACTIVATE = "ACTIVATE"
const ACTION_CODE_SUBMIT = "SUBMIT"
const ACTION_CODE_CANCEL = "CANCEL"
const ACTION_CODE_REJECT = "REJECT"
const ACTION_CODE_APPROVE = "APPROVE"

// COMMON PERMISSION
const APPROVAL_BYPASS = "appr.bypass"

// CLASS PERMISSION
const CREATE_CLASS_PERMISSION = "class.create"
const UPDATE_CLASS_PERMISSION = "class.update"
const DELETE_CLASS_PERMISSION = "class.delete"

// USER PERMISSION
const CREATE_USER_PERMISSION = "user.create"
const UPDATE_USER_PERMISSION = "user.update"
const DELETE_USER_PERMISSION = "user.delete"

// STUDENT PERMISSION
const CREATE_STUDENT_PERMISSION = "student.create"
const UPDATE_STUDENT_PERMISSION = "student.update"
const DELETE_STUDENT_PERMISSION = "student.delete"

// TEACHER AND STAFF PERMISSION
const CREATE_TNS_PERMISSION = "teacherandstaff.create"
const UPDATE_TNS_PERMISSION = "teacherandstaff.update"
const DELETE_TNS_PERMISSION = "teacherandstaff.delete"

// LOCATION PERMISSION
const LOCATION_CREATE_PERMISSION = "setting.location.create"
const LOCATION_UPDATE_PERMISSION = "setting.location.update"
const LOCATION_DELETE_PERMISSION = "setting.location.delete"

// QR CODE PERMISSION
const QR_CODE_CREATE_PERMISSION = "qrcode.create"
const QR_CODE_UPDATE_PERMISSION = "qrcode.update"
const QR_CODE_DELETE_PERMISSION = "qrcode.delete"

// const SUBMIT_CLASS_PERMISSION = "class.submit"
// const APPROVE_CLASS_PERMISSION = "class.approve"
// const REJECT_CLASS_PERMISSION = "class.reject"

// DB STATUS VAR
var DB_UUID_STATUS_ACTIVE uuid.UUID
var DB_UUID_STATUS_PENDING uuid.UUID
var DB_UUID_STATUS_INACTIVE uuid.UUID
var DB_UUID_STATUS_DELETED uuid.UUID
var DB_UUID_STATUS_NEWUSER uuid.UUID

var DB_UUID_STUDENT_ENROLLMENT_REGISTERED uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_ACTIVE uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_ON_LEAVE uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_SUSPENDED uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_INACTIVE uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_TRANSFERRED_OUT uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_WITHDRAWN uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_DROPPED_OUT uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_EXPELLED uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_GRADUATED uuid.UUID
var DB_UUID_STUDENT_ENROLLMENT_DECEASED uuid.UUID

// ROLE NAME
const ROLE_STAFF_TU = "STAFF_TU"
const ROLE_WAKA_KUR = "WAKASEK_KURIKULUM"
const ROLE_KEPSEK = "KEPALA_SEKOLAH"
const ROLE_STUDENT = "SISWA"
const ROLE_TEACHER = "GURU"
const ROLE_STAFF = "STAFF"

// STUDENT ENROLLMENT STATUS
const STUDENT_ENROLLMENT_STATUS_REGISTERED = "registered"
const STUDENT_ENROLLMENT_STATUS_ACTIVE = "active"
const STUDENT_ENROLLMENT_STATUS_ON_LEAVE = "on_leave"
const STUDENT_ENROLLMENT_STATUS_SUSPENDED = "suspended"
const STUDENT_ENROLLMENT_STATUS_INACTIVE = "inactive"
const STUDENT_ENROLLMENT_STATUS_TRANSFERRED_OUT = "transferred_out"
const STUDENT_ENROLLMENT_STATUS_WITHDRAWN = "withdrawn"
const STUDENT_ENROLLMENT_STATUS_DROPPED_OUT = "dropped_out"
const STUDENT_ENROLLMENT_STATUS_EXPELLED = "expelled"
const STUDENT_ENROLLMENT_STATUS_GRADUATED = "graduated"
const STUDENT_ENROLLMENT_STATUS_DECEASED = "deceased"

func init() {
	if API_VERSION == "" {
		API_VERSION = "v1"
	}
}

func InitVariableDB() {
	activeStatus, err := GetStatusByCode(STATUS_ACTIVE, "record_status")
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STATUS_ACTIVE = activeStatus.UUID
	}
	pendingStatus, err := GetStatusByCode(STATUS_PENDING, "record_status")
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STATUS_PENDING = pendingStatus.UUID
	}
	inactiveStatus, err := GetStatusByCode(STATUS_INACTIVE, "record_status")
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STATUS_INACTIVE = inactiveStatus.UUID
	}
	deleteStatus, err := GetStatusByCode(STATUS_DELETED, "approval")
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STATUS_DELETED = deleteStatus.UUID
	}
	newUserStatus, err := GetStatusByCode(STATUS_NEW_USER, "user_account")
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STATUS_NEWUSER = newUserStatus.UUID
	}

	log.Println("Variable Loaded")
}

func InitStudentEnrollmentStatus() {
	registeredStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_REGISTERED)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_REGISTERED = *registeredStatus
	}
	activeStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_ACTIVE)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_ACTIVE = *activeStatus
	}
	onLeaveStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_ON_LEAVE)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_ON_LEAVE = *onLeaveStatus
	}
	suspendedStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_SUSPENDED)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_SUSPENDED = *suspendedStatus
	}
	inactiveStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_INACTIVE)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_INACTIVE = *inactiveStatus
	}
	transferredOutStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_TRANSFERRED_OUT)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_TRANSFERRED_OUT = *transferredOutStatus
	}
	withdrawnStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_WITHDRAWN)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_WITHDRAWN = *withdrawnStatus
	}
	droppedOutStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_DROPPED_OUT)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_DROPPED_OUT = *droppedOutStatus
	}
	expelledStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_EXPELLED)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_EXPELLED = *expelledStatus
	}
	graduatedStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_GRADUATED)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_GRADUATED = *graduatedStatus
	}
	deceasedStatus, err := GetUserStatus(STUDENT_ENROLLMENT_STATUS_DECEASED)
	if err != nil {
		log.Fatal("Critical error loading configuration: ", err)
	} else {
		DB_UUID_STUDENT_ENROLLMENT_DECEASED = *deceasedStatus
	}

	log.Println("Student Enrollment Status Variable Loaded")
}
