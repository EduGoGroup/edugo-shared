package enum

// Permission representa un permiso del sistema RBAC en formato
// path-based (`<dominio_menu>.<recurso>.<accion>[:own]`) — D3 del
// rediseño de permisos. La gramática completa vive en
// PathPermissionRegex (permission_path.go).
type Permission string

// admin.users
const (
	PermissionUsersCreate    Permission = "admin.users.create"
	PermissionUsersRead      Permission = "admin.users.read"
	PermissionUsersUpdate    Permission = "admin.users.update"
	PermissionUsersDelete    Permission = "admin.users.delete"
	PermissionUsersReadOwn   Permission = "admin.users.read:own"
	PermissionUsersUpdateOwn Permission = "admin.users.update:own"
	// PermissionUsersGrantsManage cubre list+create+delete sobre los
	// overrides puntuales en iam.user_grants (P4-2).
	PermissionUsersGrantsManage Permission = "admin.users.grants.manage"
)

// admin.schools
const (
	PermissionSchoolsCreate Permission = "admin.schools.create"
	PermissionSchoolsRead   Permission = "admin.schools.read"
	PermissionSchoolsUpdate Permission = "admin.schools.update"
	PermissionSchoolsDelete Permission = "admin.schools.delete"
	PermissionSchoolsManage Permission = "admin.schools.manage"
)

// admin.roles
const (
	PermissionRolesCreate Permission = "admin.roles.create"
	PermissionRolesRead   Permission = "admin.roles.read"
	PermissionRolesUpdate Permission = "admin.roles.update"
	PermissionRolesDelete Permission = "admin.roles.delete"
)

// admin.permissions_mgmt
const (
	PermissionPermissionsMgmtCreate Permission = "admin.permissions_mgmt.create"
	PermissionPermissionsMgmtRead   Permission = "admin.permissions_mgmt.read"
	PermissionPermissionsMgmtUpdate Permission = "admin.permissions_mgmt.update"
	PermissionPermissionsMgmtDelete Permission = "admin.permissions_mgmt.delete"
)

// admin.screen_templates
const (
	PermissionScreenTemplatesCreate Permission = "admin.screen_templates.create"
	PermissionScreenTemplatesRead   Permission = "admin.screen_templates.read"
	PermissionScreenTemplatesUpdate Permission = "admin.screen_templates.update"
	PermissionScreenTemplatesDelete Permission = "admin.screen_templates.delete"
)

// admin.screen_instances
const (
	PermissionScreenInstancesCreate Permission = "admin.screen_instances.create"
	PermissionScreenInstancesRead   Permission = "admin.screen_instances.read"
	PermissionScreenInstancesUpdate Permission = "admin.screen_instances.update"
	PermissionScreenInstancesDelete Permission = "admin.screen_instances.delete"
)

// admin.audit
const (
	PermissionAuditRead   Permission = "admin.audit.read"
	PermissionAuditExport Permission = "admin.audit.export"
)

// admin.concept_types
const (
	PermissionConceptTypesCreate Permission = "admin.concept_types.create"
	PermissionConceptTypesRead   Permission = "admin.concept_types.read"
	PermissionConceptTypesUpdate Permission = "admin.concept_types.update"
	PermissionConceptTypesDelete Permission = "admin.concept_types.delete"
)

// admin.system_settings
const (
	PermissionSystemSettingsSettings Permission = "admin.system_settings.settings"
	PermissionSystemSettingsRead     Permission = "admin.system_settings.read"
	PermissionSystemSettingsUpdate   Permission = "admin.system_settings.update"
)

// academic.units
const (
	PermissionUnitsCreate Permission = "academic.units.create"
	PermissionUnitsRead   Permission = "academic.units.read"
	PermissionUnitsUpdate Permission = "academic.units.update"
	PermissionUnitsDelete Permission = "academic.units.delete"
)

// academic.memberships
const (
	PermissionMembershipsCreate Permission = "academic.memberships.create"
	PermissionMembershipsRead   Permission = "academic.memberships.read"
	PermissionMembershipsUpdate Permission = "academic.memberships.update"
	PermissionMembershipsDelete Permission = "academic.memberships.delete"
)

// academic.my_memberships
const (
	// PermissionMyMembershipsReadOwn permite al alumno leer SOLO sus
	// propias membresías ("mis materias"), sin listar las de otros ni la
	// unidad completa. Lo usa el rol student. El self-check vive en el
	// handler GET /users/:user_id/memberships (plan 006 N1.C). Vive bajo
	// un path propio (academic.my_memberships.*) para que el gate de menú
	// por path-prefix NO haga aparecer el item admin "memberships".
	PermissionMyMembershipsReadOwn Permission = "academic.my_memberships.read:own"
)

// academic.subjects
const (
	PermissionSubjectsCreate Permission = "academic.subjects.create"
	PermissionSubjectsRead   Permission = "academic.subjects.read"
	PermissionSubjectsUpdate Permission = "academic.subjects.update"
	PermissionSubjectsDelete Permission = "academic.subjects.delete"
)

// academic.subject_offerings (sesiones de materia, ADR 0009 / plan 010 N1.7)
const (
	PermissionSubjectOfferingsCreate Permission = "academic.subject_offerings.create"
	PermissionSubjectOfferingsRead   Permission = "academic.subject_offerings.read"
	PermissionSubjectOfferingsUpdate Permission = "academic.subject_offerings.update"
	PermissionSubjectOfferingsDelete Permission = "academic.subject_offerings.delete"
	// PermissionSubjectOfferingsEnroll cubre alta y baja de matrícula
	// (inscripción por lote a una sesión).
	PermissionSubjectOfferingsEnroll Permission = "academic.subject_offerings.enroll"
)

// academic.guardian_relations
const (
	PermissionGuardianRelationsRead    Permission = "academic.guardian_relations.read"
	PermissionGuardianRelationsApprove Permission = "academic.guardian_relations.approve"
	PermissionGuardianRelationsRequest Permission = "academic.guardian_relations.request"
	PermissionGuardianRelationsManage  Permission = "academic.guardian_relations.manage"
)

// academic.invitations
const (
	PermissionInvitationsCreate Permission = "academic.invitations.create"
	PermissionInvitationsRead   Permission = "academic.invitations.read"
	PermissionInvitationsRevoke Permission = "academic.invitations.revoke"
)

// academic.join_requests
const (
	PermissionJoinRequestsRead   Permission = "academic.join_requests.read"
	PermissionJoinRequestsReject Permission = "academic.join_requests.reject"
)

// academic.periods
const (
	PermissionPeriodsCreate   Permission = "academic.periods.create"
	PermissionPeriodsRead     Permission = "academic.periods.read"
	PermissionPeriodsUpdate   Permission = "academic.periods.update"
	PermissionPeriodsDelete   Permission = "academic.periods.delete"
	PermissionPeriodsActivate Permission = "academic.periods.activate"
)

// academic.grades
const (
	PermissionGradesCreate   Permission = "academic.grades.create"
	PermissionGradesRead     Permission = "academic.grades.read"
	PermissionGradesUpdate   Permission = "academic.grades.update"
	PermissionGradesFinalize Permission = "academic.grades.finalize"
)

// academic.attendance
const (
	PermissionAttendanceCreate Permission = "academic.attendance.create"
	PermissionAttendanceRead   Permission = "academic.attendance.read"
	PermissionAttendanceUpdate Permission = "academic.attendance.update"
)

// academic.schedules
const (
	PermissionSchedulesCreate Permission = "academic.schedules.create"
	PermissionSchedulesRead   Permission = "academic.schedules.read"
	PermissionSchedulesUpdate Permission = "academic.schedules.update"
	PermissionSchedulesDelete Permission = "academic.schedules.delete"
)

// academic.calendar
const (
	PermissionCalendarEventsCreate Permission = "academic.calendar.create"
	PermissionCalendarEventsRead   Permission = "academic.calendar.read"
	PermissionCalendarEventsUpdate Permission = "academic.calendar.update"
	PermissionCalendarEventsDelete Permission = "academic.calendar.delete"
)

// academic.announcements
const (
	PermissionAnnouncementsCreate Permission = "academic.announcements.create"
	PermissionAnnouncementsRead   Permission = "academic.announcements.read"
	PermissionAnnouncementsUpdate Permission = "academic.announcements.update"
	PermissionAnnouncementsDelete Permission = "academic.announcements.delete"
)

// platform.colors
// Recurso CRUD plano demo introducido en la Fase 3 SDUI (F3-REQ-4) para
// validar que crear un CRUD nuevo en EduGo NO requiere código Kotlin.
const (
	PermissionColorsCreate Permission = "platform.colors.create"
	PermissionColorsRead   Permission = "platform.colors.read"
	PermissionColorsUpdate Permission = "platform.colors.update"
	PermissionColorsDelete Permission = "platform.colors.delete"
)

// content.materials
const (
	PermissionMaterialsCreate   Permission = "content.materials.create"
	PermissionMaterialsRead     Permission = "content.materials.read"
	PermissionMaterialsUpdate   Permission = "content.materials.update"
	PermissionMaterialsDelete   Permission = "content.materials.delete"
	PermissionMaterialsPublish  Permission = "content.materials.publish"
	PermissionMaterialsDownload Permission = "content.materials.download"
	PermissionMaterialsUpload   Permission = "content.materials.upload"
)

// content.assessments
const (
	PermissionAssessmentsCreate      Permission = "content.assessments.create"
	PermissionAssessmentsRead        Permission = "content.assessments.read"
	PermissionAssessmentsUpdate      Permission = "content.assessments.update"
	PermissionAssessmentsDelete      Permission = "content.assessments.delete"
	PermissionAssessmentsPublish     Permission = "content.assessments.publish"
	PermissionAssessmentsGrade       Permission = "content.assessments.grade"
	PermissionAssessmentsAttempt     Permission = "content.assessments.attempt"
	PermissionAssessmentsViewResults Permission = "content.assessments.view_results"
	PermissionAssessmentsAssign      Permission = "content.assessments.assign"
	PermissionAssessmentsReview      Permission = "content.assessments.review"
)

// content.assessments_student
const (
	PermissionAssessmentsStudentRead Permission = "content.assessments_student.read"
)

// reports.progress
const (
	PermissionProgressRead    Permission = "reports.progress.read"
	PermissionProgressUpdate  Permission = "reports.progress.update"
	PermissionProgressReadOwn Permission = "reports.progress.read:own"
)

// reports.stats
const (
	PermissionStatsGlobal Permission = "reports.stats.global"
	PermissionStatsSchool Permission = "reports.stats.school"
	PermissionStatsUnit   Permission = "reports.stats.unit"
)

// roots de 2 segmentos (recursos sin parent_id)
const (
	PermissionDashboardView        Permission = "dashboard.view"
	PermissionMenuRead             Permission = "menu.read"
	PermissionMenuFullRead         Permission = "menu.full_read"
	PermissionNotificationsRead    Permission = "notifications.read"
	PermissionScreensRead          Permission = "screens.read"
	PermissionContextBrowseSchools Permission = "context.browse_schools"
	PermissionContextBrowseUnits   Permission = "context.browse_units"
	PermissionReportsRead          Permission = "reports.read"
)

// String retorna la representación en string del permiso.
func (p Permission) String() string {
	return string(p)
}

// IsValid verifica si el permiso es uno de los conocidos por el sistema
// (formato exacto, no patterns con wildcards).
func (p Permission) IsValid() bool {
	return AllPermissions[p]
}

// AllPermissions es el catálogo cerrado de permisos conocidos del
// sistema. Cualquier `Permission*` declarado arriba debe aparecer acá;
// el test `TestAllPermissions_MapIntegrity` lo verifica.
var AllPermissions = map[Permission]bool{
	// admin.users
	PermissionUsersCreate:       true,
	PermissionUsersRead:         true,
	PermissionUsersUpdate:       true,
	PermissionUsersDelete:       true,
	PermissionUsersReadOwn:      true,
	PermissionUsersUpdateOwn:    true,
	PermissionUsersGrantsManage: true,
	// admin.schools
	PermissionSchoolsCreate: true,
	PermissionSchoolsRead:   true,
	PermissionSchoolsUpdate: true,
	PermissionSchoolsDelete: true,
	PermissionSchoolsManage: true,
	// admin.roles
	PermissionRolesCreate: true,
	PermissionRolesRead:   true,
	PermissionRolesUpdate: true,
	PermissionRolesDelete: true,
	// admin.permissions_mgmt
	PermissionPermissionsMgmtCreate: true,
	PermissionPermissionsMgmtRead:   true,
	PermissionPermissionsMgmtUpdate: true,
	PermissionPermissionsMgmtDelete: true,
	// admin.screen_templates
	PermissionScreenTemplatesCreate: true,
	PermissionScreenTemplatesRead:   true,
	PermissionScreenTemplatesUpdate: true,
	PermissionScreenTemplatesDelete: true,
	// admin.screen_instances
	PermissionScreenInstancesCreate: true,
	PermissionScreenInstancesRead:   true,
	PermissionScreenInstancesUpdate: true,
	PermissionScreenInstancesDelete: true,
	// admin.audit
	PermissionAuditRead:   true,
	PermissionAuditExport: true,
	// admin.concept_types
	PermissionConceptTypesCreate: true,
	PermissionConceptTypesRead:   true,
	PermissionConceptTypesUpdate: true,
	PermissionConceptTypesDelete: true,
	// admin.system_settings
	PermissionSystemSettingsSettings: true,
	PermissionSystemSettingsRead:     true,
	PermissionSystemSettingsUpdate:   true,
	// academic.units
	PermissionUnitsCreate: true,
	PermissionUnitsRead:   true,
	PermissionUnitsUpdate: true,
	PermissionUnitsDelete: true,
	// academic.memberships
	PermissionMembershipsCreate:  true,
	PermissionMembershipsRead:    true,
	PermissionMembershipsUpdate:  true,
	PermissionMembershipsDelete:  true,
	// academic.my_memberships
	PermissionMyMembershipsReadOwn: true,
	// academic.subjects
	PermissionSubjectsCreate: true,
	PermissionSubjectsRead:   true,
	PermissionSubjectsUpdate: true,
	PermissionSubjectsDelete: true,
	// academic.subject_offerings
	PermissionSubjectOfferingsCreate: true,
	PermissionSubjectOfferingsRead:   true,
	PermissionSubjectOfferingsUpdate: true,
	PermissionSubjectOfferingsDelete: true,
	PermissionSubjectOfferingsEnroll: true,
	// academic.guardian_relations
	PermissionGuardianRelationsRead:    true,
	PermissionGuardianRelationsApprove: true,
	PermissionGuardianRelationsRequest: true,
	PermissionGuardianRelationsManage:  true,
	// academic.invitations
	PermissionInvitationsCreate: true,
	PermissionInvitationsRead:   true,
	PermissionInvitationsRevoke: true,
	// academic.join_requests
	PermissionJoinRequestsRead:   true,
	PermissionJoinRequestsReject: true,
	// academic.periods
	PermissionPeriodsCreate:   true,
	PermissionPeriodsRead:     true,
	PermissionPeriodsUpdate:   true,
	PermissionPeriodsDelete:   true,
	PermissionPeriodsActivate: true,
	// academic.grades
	PermissionGradesCreate:   true,
	PermissionGradesRead:     true,
	PermissionGradesUpdate:   true,
	PermissionGradesFinalize: true,
	// academic.attendance
	PermissionAttendanceCreate: true,
	PermissionAttendanceRead:   true,
	PermissionAttendanceUpdate: true,
	// academic.schedules
	PermissionSchedulesCreate: true,
	PermissionSchedulesRead:   true,
	PermissionSchedulesUpdate: true,
	PermissionSchedulesDelete: true,
	// academic.calendar
	PermissionCalendarEventsCreate: true,
	PermissionCalendarEventsRead:   true,
	PermissionCalendarEventsUpdate: true,
	PermissionCalendarEventsDelete: true,
	// academic.announcements
	PermissionAnnouncementsCreate: true,
	PermissionAnnouncementsRead:   true,
	PermissionAnnouncementsUpdate: true,
	PermissionAnnouncementsDelete: true,
	// platform.colors
	PermissionColorsCreate: true,
	PermissionColorsRead:   true,
	PermissionColorsUpdate: true,
	PermissionColorsDelete: true,
	// content.materials
	PermissionMaterialsCreate:   true,
	PermissionMaterialsRead:     true,
	PermissionMaterialsUpdate:   true,
	PermissionMaterialsDelete:   true,
	PermissionMaterialsPublish:  true,
	PermissionMaterialsDownload: true,
	PermissionMaterialsUpload:   true,
	// content.assessments
	PermissionAssessmentsCreate:      true,
	PermissionAssessmentsRead:        true,
	PermissionAssessmentsUpdate:      true,
	PermissionAssessmentsDelete:      true,
	PermissionAssessmentsPublish:     true,
	PermissionAssessmentsGrade:       true,
	PermissionAssessmentsAttempt:     true,
	PermissionAssessmentsViewResults: true,
	PermissionAssessmentsAssign:      true,
	PermissionAssessmentsReview:      true,
	// content.assessments_student
	PermissionAssessmentsStudentRead: true,
	// reports.progress
	PermissionProgressRead:    true,
	PermissionProgressUpdate:  true,
	PermissionProgressReadOwn: true,
	// reports.stats
	PermissionStatsGlobal: true,
	PermissionStatsSchool: true,
	PermissionStatsUnit:   true,
	// roots de 2 segmentos
	PermissionDashboardView:        true,
	PermissionMenuRead:             true,
	PermissionMenuFullRead:         true,
	PermissionNotificationsRead:    true,
	PermissionScreensRead:          true,
	PermissionContextBrowseSchools: true,
	PermissionContextBrowseUnits:   true,
	PermissionReportsRead:          true,
}
