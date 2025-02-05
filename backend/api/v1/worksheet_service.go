package v1

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	"github.com/bytebase/bytebase/backend/store"
	"github.com/bytebase/bytebase/backend/utils"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

// WorksheetService implements the worksheet service.
type WorksheetService struct {
	v1pb.UnimplementedWorksheetServiceServer
	store *store.Store
}

// NewWorksheetService creates a new WorksheetService.
func NewWorksheetService(store *store.Store) *WorksheetService {
	return &WorksheetService{
		store: store,
	}
}

// CreateWorksheet creates a new worksheet.
func (s *WorksheetService) CreateWorksheet(ctx context.Context, request *v1pb.CreateWorksheetRequest) (*v1pb.Worksheet, error) {
	if request.Worksheet == nil {
		return nil, status.Errorf(codes.InvalidArgument, "worksheet must be set")
	}
	principalID, ok := ctx.Value(common.PrincipalIDContextKey).(int)
	if !ok {
		return nil, status.Errorf(codes.Internal, "principal ID not found")
	}

	projectResourceID, err := common.GetProjectID(request.Worksheet.Project)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		ResourceID: &projectResourceID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get project with resource id %q, err: %s", projectResourceID, err.Error()))
	}
	if project == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q not found", projectResourceID))
	}
	if project.Deleted {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q had deleted", projectResourceID))
	}

	var databaseUID *int
	if request.Worksheet.Database != "" {
		instanceResourceID, databaseName, err := common.GetInstanceDatabaseID(request.Worksheet.Database)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		instance, err := s.store.GetInstanceV2(ctx, &store.FindInstanceMessage{
			ResourceID: &instanceResourceID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get instance with resource id %q, err: %s", instanceResourceID, err.Error()))
		}
		if instance == nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("instance with resource id %q not found", instanceResourceID))
		}

		find := &store.FindDatabaseMessage{
			ProjectID:  &projectResourceID,
			InstanceID: &instanceResourceID,
		}
		// It's chaos. We return /instance/{resource id}/databases/{uid} database in find worksheet request,
		// but the frontend use both /instance/{resource id}/databases/{uid} and /instance/{resource id}/databases/{name}, sometimes the name will convert to int id incorrectly.
		// For database v1 api, we should only use the /instance/{resource id}/databases/{name}
		// We need to remove legacy code after the migration.
		dbUID, isNumber := isNumber(databaseName)
		if instanceResourceID == "-" && isNumber {
			find.UID = &dbUID
		} else {
			find.DatabaseName = &databaseName
			find.IgnoreCaseSensitive = store.IgnoreDatabaseAndTableCaseSensitive(instance)
		}

		database, err := s.store.GetDatabaseV2(ctx, find)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get database with name %q, err: %s", databaseName, err.Error()))
		}
		if database == nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("database with name %q not found in project %q instance %q", databaseName, projectResourceID, instanceResourceID))
		}
		databaseUID = &database.UID
	}
	storeWorksheetCreate, err := convertToStoreWorksheetMessage(project.UID, databaseUID, principalID, request.Worksheet)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("failed to convert worksheet: %v", err))
	}
	worksheet, err := s.store.CreateWorkSheet(ctx, storeWorksheetCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to create worksheet: %v", err))
	}
	v1pbWorksheet, err := s.convertToAPIWorksheetMessage(ctx, worksheet)
	if err != nil {
		return nil, err
	}
	return v1pbWorksheet, nil
}

// GetWorksheet returns the requested worksheet, cutoff the content if the content is too long and the `raw` flag in request is false.
func (s *WorksheetService) GetWorksheet(ctx context.Context, request *v1pb.GetWorksheetRequest) (*v1pb.Worksheet, error) {
	worksheetUID, err := common.GetWorksheetUID(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if worksheetUID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid worksheet id %d, must be positive integer", worksheetUID))
	}

	find := &store.FindWorkSheetMessage{
		UID:      &worksheetUID,
		LoadFull: true,
	}
	worksheet, err := s.findWorksheet(ctx, find)
	if err != nil {
		return nil, err
	}

	canAccess, err := s.canReadWorksheet(ctx, worksheet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to check access with error: %v", err))
	}
	if !canAccess {
		return nil, status.Errorf(codes.PermissionDenied, "cannot access worksheet %s", worksheet.Title)
	}

	v1pbWorksheet, err := s.convertToAPIWorksheetMessage(ctx, worksheet)
	if err != nil {
		return nil, err
	}
	return v1pbWorksheet, nil
}

// SearchWorksheets returns a list of worksheets based on the search filters.
func (s *WorksheetService) SearchWorksheets(ctx context.Context, request *v1pb.SearchWorksheetsRequest) (*v1pb.SearchWorksheetsResponse, error) {
	principalID, ok := ctx.Value(common.PrincipalIDContextKey).(int)
	if !ok {
		return nil, status.Errorf(codes.Internal, "principal ID not found")
	}

	worksheetFind := &store.FindWorkSheetMessage{}
	// TODO(zp): It is difficult to find all the worksheets visible to a principal atomically
	// without adding a new store layer method, which has two parts:
	// 1. creator = principal && visibility in (PROJECT, PUBLIC, PRIVATE)
	// 2. creator ! = principal && visibility in (PROJECT, PUBLIC)
	// So we don't allow empty filter for now.
	if request.Filter == "" {
		return nil, status.Errorf(codes.InvalidArgument, "filter should not be empty")
	}

	specs, err := parseFilter(request.Filter)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	for _, spec := range specs {
		switch spec.key {
		case "creator":
			creatorEmail := strings.TrimPrefix(spec.value, "users/")
			if creatorEmail == "" {
				return nil, status.Errorf(codes.InvalidArgument, "invalid empty creator identifier")
			}
			user, err := s.store.GetUser(ctx, &store.FindUserMessage{
				Email: &creatorEmail,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get user: %s", err.Error()))
			}
			if user == nil {
				return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with email %s not found", creatorEmail))
			}
			switch spec.operator {
			case comparatorTypeEqual:
				worksheetFind.CreatorID = &user.ID
				worksheetFind.Visibilities = []store.WorkSheetVisibility{store.ProjectWorkSheet, store.PublicWorkSheet, store.PrivateWorkSheet}
			case comparatorTypeNotEqual:
				worksheetFind.ExcludedCreatorID = &user.ID
				worksheetFind.Visibilities = []store.WorkSheetVisibility{store.ProjectWorkSheet, store.PublicWorkSheet}
				worksheetFind.PrincipalID = &user.ID
			default:
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid operator %q for creator", spec.operator))
			}
		case "starred":
			if spec.operator != comparatorTypeEqual {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid operator %q for starred", spec.operator))
			}
			switch spec.value {
			case "true":
				worksheetFind.OrganizerPrincipalIDStarred = &principalID
			case "false":
				worksheetFind.OrganizerPrincipalIDNotStarred = &principalID
			default:
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid value %q for starred", spec.value))
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid filter key %q", spec.key))
		}
	}
	worksheetList, err := s.store.ListWorkSheets(ctx, worksheetFind, principalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to list worksheets: %v", err))
	}

	var v1pbWorksheets []*v1pb.Worksheet
	for _, worksheet := range worksheetList {
		canAccess, err := s.canReadWorksheet(ctx, worksheet)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to check access with error: %v", err))
		}
		if !canAccess {
			slog.Warn("cannot access worksheet", slog.String("name", worksheet.Title))
			continue
		}
		v1pbWorksheet, err := s.convertToAPIWorksheetMessage(ctx, worksheet)
		if err != nil {
			st := status.Convert(err)
			if st.Code() == codes.NotFound {
				slog.Debug("failed to found resource for worksheet", log.BBError(err), slog.Int("id", worksheet.UID), slog.Int("project", worksheet.ProjectUID))
				continue
			}
			return nil, err
		}
		v1pbWorksheets = append(v1pbWorksheets, v1pbWorksheet)
	}
	return &v1pb.SearchWorksheetsResponse{
		Worksheets: v1pbWorksheets,
	}, nil
}

// UpdateWorksheet updates a worksheet.
func (s *WorksheetService) UpdateWorksheet(ctx context.Context, request *v1pb.UpdateWorksheetRequest) (*v1pb.Worksheet, error) {
	if request.Worksheet == nil {
		return nil, status.Errorf(codes.InvalidArgument, "worksheet cannot be empty")
	}
	if request.UpdateMask == nil {
		return nil, status.Errorf(codes.InvalidArgument, "update mask cannot be empty")
	}
	if request.Worksheet.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "worksheet name cannot be empty")
	}

	worksheetUID, err := common.GetWorksheetUID(request.Worksheet.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if worksheetUID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid worksheet id %d, must be positive integer", worksheetUID))
	}

	principalID, ok := ctx.Value(common.PrincipalIDContextKey).(int)
	if !ok {
		return nil, status.Errorf(codes.Internal, "principal ID not found")
	}
	worksheet, err := s.store.GetWorkSheet(ctx, &store.FindWorkSheetMessage{
		UID: &worksheetUID,
	}, principalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get worksheet: %v", err))
	}
	if worksheet == nil {
		return nil, status.Errorf(codes.NotFound, "worksheet %q not found", request.Worksheet.Name)
	}
	canAccess, err := s.canWriteWorksheet(ctx, worksheet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to check access with error: %v", err))
	}
	if !canAccess {
		return nil, status.Errorf(codes.PermissionDenied, "cannot write worksheet %s", worksheet.Title)
	}

	worksheetPatch := &store.PatchWorkSheetMessage{
		UID:       worksheet.UID,
		UpdaterID: principalID,
	}
	for _, path := range request.UpdateMask.Paths {
		switch path {
		case "title":
			worksheetPatch.Title = &request.Worksheet.Title
		case "content":
			statement := string(request.Worksheet.Content)
			worksheetPatch.Statement = &statement
		case "visibility":
			visibility, err := convertToStoreWorksheetVisibility(request.Worksheet.Visibility)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid visibility %q", request.Worksheet.Visibility))
			}
			stringVisibility := string(visibility)
			worksheetPatch.Visibility = &stringVisibility
		default:
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid update mask path %q", path))
		}
	}
	if err := s.store.PatchWorkSheet(ctx, worksheetPatch); err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to update worksheet: %v", err))
	}

	worksheet, err = s.store.GetWorkSheet(ctx, &store.FindWorkSheetMessage{
		UID: &worksheetUID,
	}, principalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get worksheet: %v", err))
	}
	if worksheet == nil {
		return nil, status.Errorf(codes.NotFound, "worksheet %q not found", request.Worksheet.Name)
	}
	v1pbWorksheet, err := s.convertToAPIWorksheetMessage(ctx, worksheet)
	if err != nil {
		return nil, err
	}

	return v1pbWorksheet, nil
}

// DeleteWorksheet deletes a worksheet.
func (s *WorksheetService) DeleteWorksheet(ctx context.Context, request *v1pb.DeleteWorksheetRequest) (*emptypb.Empty, error) {
	worksheetUID, err := common.GetWorksheetUID(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	principalID, ok := ctx.Value(common.PrincipalIDContextKey).(int)
	if !ok {
		return nil, status.Errorf(codes.Internal, "principal ID not found")
	}
	worksheet, err := s.store.GetWorkSheet(ctx, &store.FindWorkSheetMessage{
		UID: &worksheetUID,
	}, principalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get worksheet: %v", err))
	}
	if worksheet == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("worksheet with id %d not found", worksheetUID))
	}
	canAccess, err := s.canWriteWorksheet(ctx, worksheet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to check access with error: %v", err))
	}
	if !canAccess {
		return nil, status.Errorf(codes.PermissionDenied, "cannot write worksheet %s", worksheet.Title)
	}

	if err := s.store.DeleteWorkSheet(ctx, worksheetUID); err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to delete worksheet: %v", err))
	}

	return &emptypb.Empty{}, nil
}

// UpdateWorksheetOrganizer upsert the worksheet organizer.
func (s *WorksheetService) UpdateWorksheetOrganizer(ctx context.Context, request *v1pb.UpdateWorksheetOrganizerRequest) (*v1pb.WorksheetOrganizer, error) {
	worksheetUID, err := common.GetWorksheetUID(request.Organizer.Worksheet)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if worksheetUID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid worksheet id %d, must be positive integer", worksheetUID))
	}

	worksheet, err := s.findWorksheet(ctx, &store.FindWorkSheetMessage{
		UID: &worksheetUID,
	})
	if err != nil {
		return nil, err
	}

	canAccess, err := s.canReadWorksheet(ctx, worksheet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to check access with error: %v", err))
	}
	if !canAccess {
		return nil, status.Errorf(codes.PermissionDenied, "cannot access worksheet %s", worksheet.Title)
	}

	principalID, ok := ctx.Value(common.PrincipalIDContextKey).(int)
	if !ok {
		return nil, status.Errorf(codes.Internal, "principal ID not found")
	}
	worksheetOrganizerUpsert := &store.WorksheetOrganizerMessage{
		WorksheetUID: worksheetUID,
		PrincipalUID: principalID,
	}

	for _, path := range request.UpdateMask.Paths {
		if path == "starred" {
			worksheetOrganizerUpsert.Starred = request.Organizer.Starred
		}
	}

	organizer, err := s.store.UpsertWorksheetOrganizer(ctx, worksheetOrganizerUpsert)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upsert organizer for worksheet %s with error: %v", request.Organizer.Worksheet, err)
	}

	return &v1pb.WorksheetOrganizer{
		Worksheet: request.Organizer.Worksheet,
		Starred:   organizer.Starred,
	}, nil
}

func (s *WorksheetService) findWorksheet(ctx context.Context, find *store.FindWorkSheetMessage) (*store.WorkSheetMessage, error) {
	principalID, ok := ctx.Value(common.PrincipalIDContextKey).(int)
	if !ok {
		return nil, status.Errorf(codes.Internal, "principal ID not found")
	}
	worksheet, err := s.store.GetWorkSheet(ctx, find, principalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get worksheet: %v", err))
	}
	if worksheet == nil {
		return nil, status.Errorf(codes.NotFound, "cannot find the worksheet")
	}
	return worksheet, nil
}

// canWriteWorksheet check if the principal can write the worksheet.
// worksheet if writable when:
// PRIVATE: the creator only.
// PROJECT: the creator or project role can manage worksheet, workspace Owner and DBA.
// PUBLIC: the creator only.
func (s *WorksheetService) canWriteWorksheet(ctx context.Context, worksheet *store.WorkSheetMessage) (bool, error) {
	user, ok := ctx.Value(common.UserContextKey).(*store.UserMessage)
	if !ok {
		return false, status.Errorf(codes.Internal, "user not found")
	}

	if worksheet.CreatorID == user.ID {
		return true, nil
	}

	if worksheet.Visibility == store.ProjectWorkSheet {
		projectRoles, err := s.findProjectRoles(ctx, worksheet.ProjectUID, user)
		if err != nil {
			return false, err
		}
		if len(projectRoles) == 0 {
			return false, nil
		}
		return projectRoles[api.ProjectOwner], nil
	}

	return false, nil
}

// canReadWorksheet check if the principal can read the worksheet.
// worksheet is readable when:
// PRIVATE: the creator only.
// PROJECT: the creator and members in the project.
// PUBLIC: everyone in the workspace.
func (s *WorksheetService) canReadWorksheet(ctx context.Context, worksheet *store.WorkSheetMessage) (bool, error) {
	user, ok := ctx.Value(common.UserContextKey).(*store.UserMessage)
	if !ok {
		return false, status.Errorf(codes.Internal, "user not found")
	}

	switch worksheet.Visibility {
	case store.PrivateWorkSheet:
		return worksheet.CreatorID == user.ID, nil
	case store.PublicWorkSheet:
		return true, nil
	case store.ProjectWorkSheet:
		if slices.Contains(user.Roles, api.WorkspaceAdmin) || slices.Contains(user.Roles, api.WorkspaceDBA) {
			return true, nil
		}
		projectRoles, err := s.findProjectRoles(ctx, worksheet.ProjectUID, user)
		if err != nil {
			return false, err
		}
		return len(projectRoles) > 0, nil
	}
	return false, nil
}

func (s *WorksheetService) findProjectRoles(ctx context.Context, projectUID int, user *store.UserMessage) (map[api.Role]bool, error) {
	policy, err := s.store.GetProjectPolicy(ctx, &store.GetProjectPolicyMessage{UID: &projectUID})
	if err != nil {
		return nil, err
	}
	return utils.GetUserRolesMap(user, policy)
}

func (s *WorksheetService) convertToAPIWorksheetMessage(ctx context.Context, worksheet *store.WorkSheetMessage) (*v1pb.Worksheet, error) {
	databaseParent := ""
	if worksheet.DatabaseUID != nil {
		database, err := s.store.GetDatabaseV2(ctx, &store.FindDatabaseMessage{
			UID: worksheet.DatabaseUID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get database: %v", err))
		}
		if database == nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("database with id %d not found", *worksheet.DatabaseUID))
		}
		databaseParent = fmt.Sprintf("%s%s/%s%s", common.InstanceNamePrefix, database.InstanceID, common.DatabaseIDPrefix, database.DatabaseName)
	}

	visibility := v1pb.Worksheet_VISIBILITY_UNSPECIFIED
	switch worksheet.Visibility {
	case store.PublicWorkSheet:
		visibility = v1pb.Worksheet_VISIBILITY_PUBLIC
	case store.ProjectWorkSheet:
		visibility = v1pb.Worksheet_VISIBILITY_PROJECT
	case store.PrivateWorkSheet:
		visibility = v1pb.Worksheet_VISIBILITY_PRIVATE
	}

	creator, err := s.store.GetUserByID(ctx, worksheet.CreatorID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get creator: %v", err))
	}

	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		UID: &worksheet.ProjectUID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get project: %v", err))
	}
	if project == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with id %d not found", worksheet.ProjectUID))
	}
	return &v1pb.Worksheet{
		Name:        fmt.Sprintf("%s%d", common.WorksheetIDPrefix, worksheet.UID),
		Project:     fmt.Sprintf("%s%s", common.ProjectNamePrefix, project.ResourceID),
		Database:    databaseParent,
		Title:       worksheet.Title,
		Creator:     fmt.Sprintf("users/%s", creator.Email),
		CreateTime:  timestamppb.New(worksheet.CreatedTime),
		UpdateTime:  timestamppb.New(worksheet.UpdatedTime),
		Content:     []byte(worksheet.Statement),
		ContentSize: worksheet.Size,
		Visibility:  visibility,
		Starred:     worksheet.Starred,
	}, nil
}

func convertToStoreWorksheetMessage(projectUID int, databaseUID *int, creatorID int, worksheet *v1pb.Worksheet) (*store.WorkSheetMessage, error) {
	visibility, err := convertToStoreWorksheetVisibility(worksheet.Visibility)
	if err != nil {
		return nil, err
	}

	worksheetMessage := &store.WorkSheetMessage{
		ProjectUID:  projectUID,
		DatabaseUID: databaseUID,
		CreatorID:   creatorID,
		Title:       worksheet.Title,
		Statement:   string(worksheet.Content),
		Visibility:  visibility,
	}

	return worksheetMessage, nil
}

func convertToStoreWorksheetVisibility(visibility v1pb.Worksheet_Visibility) (store.WorkSheetVisibility, error) {
	switch visibility {
	case v1pb.Worksheet_VISIBILITY_UNSPECIFIED:
		return store.WorkSheetVisibility(""), status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid visibility %q", visibility))
	case v1pb.Worksheet_VISIBILITY_PUBLIC:
		return store.PublicWorkSheet, nil
	case v1pb.Worksheet_VISIBILITY_PROJECT:
		return store.ProjectWorkSheet, nil
	case v1pb.Worksheet_VISIBILITY_PRIVATE:
		return store.PrivateWorkSheet, nil
	default:
		return store.WorkSheetVisibility(""), status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid visibility %q", visibility))
	}
}
