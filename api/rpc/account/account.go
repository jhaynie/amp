package account

import (
	"fmt"
	pb "github.com/golang/protobuf/ptypes/empty"
	"github.com/hlandau/passlib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const hash = "$s2$16384$8$1$42JtddBgSqrJMwc3YuTNW+R+$ISfEF3jkvYQYk4AK/UFAxdqnmNFVeUw2gUVXEMBDAng=" // password

// Server is used to implement account.AccountServer
type Server struct{}

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (out *SignUpReply, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	_, err = passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "hashing error")
	}
	out = &SignUpReply{}
	out.SessionKey = in.Name
	return
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (out *pb.Empty, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	_, err = passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "hashing error")
	}
	out = &pb.Empty{}
	fmt.Println(in.Code)
	return
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (out *pb.Empty, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	out = &pb.Empty{}
	return
}

// PasswordChange implements account.PasswordChange
func (s *Server) PasswordChange(ctx context.Context, in *PasswordChangeRequest) (out *pb.Empty, err error) {
	return
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (out *LogInReply, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	_, err = passlib.Verify(in.Password, hash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}
	out = &LogInReply{}
	out.SessionKey = in.Name
	return
}

// List implements account.List
func (s *Server) List(ctx context.Context, in *ListAccountRequest) (out *ListAccountReply, err error) {
	//if in.Type != "individual" && in.Type != "organization" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "account type is mandatory")
	//}
	out = &ListAccountReply{}
	return
}

// Switch implements account.Switch
func (s *Server) Switch(ctx context.Context, in *SwitchRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &pb.Empty{}
	return
}

// GetDetails implements account.GetDetails
func (s *Server) GetDetails(ctx context.Context, in *GetAccountDetailsRequest) (out *GetAccountDetailsReply, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	out = &GetAccountDetailsReply{}
	return
}

// Edit implements account.Edit
func (s *Server) Edit(ctx context.Context, in *EditAccountRequest) (out *pb.Empty, err error) {
	err = in.Validate()
	if err != nil {
		return
	}
	//if in.NewPassword != "" {
	//	_, err := passlib.Verify(in.Password, hash)
	//	if err != nil {
	//		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	//	}
	//	_, err = passlib.Hash(in.NewPassword)
	//	if err != nil {
	//		return nil, grpc.Errorf(codes.Internal, "hashing error")
	//	}
	//}
	out = &pb.Empty{}
	return
}

// Delete implements account.Delete
func (s *Server) Delete(ctx context.Context, in *DeleteAccountRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	out = &pb.Empty{}
	return
}

// GetTeams implements account.GetTeams
func (s *Server) GetTeams(ctx context.Context, in *GetTeamsRequest) (out *GetTeamsReply, err error) {
	out = &GetTeamsReply{}
	return
}

// AddOrganizationMemberships implements account.AddOrganizationMemberships
func (s *Server) AddOrganizationMemberships(ctx context.Context, in *AddOrganizationMembershipsRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	out = &pb.Empty{}
	return
}

// DeleteOrganizationMemberships implements account.DeleteOrganizationMemberships
func (s *Server) DeleteOrganizationMemberships(ctx context.Context, in *DeleteOrganizationMembershipsRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	out = &pb.Empty{}
	return
}

// CreateTeam implements account.CreateTeam
func (s *Server) CreateTeam(ctx context.Context, in *CreateTeamRequest) (out *pb.Empty, err error) {
	out = &pb.Empty{}
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	return
}

// ListTeam implements account.ListTeam
func (s *Server) ListTeam(ctx context.Context, in *ListTeamRequest) (out *ListTeamReply, err error) {
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &ListTeamReply{}
	return
}

// EditTeam implements account.EditTeam
func (s *Server) EditTeam(ctx context.Context, in *EditTeamRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &pb.Empty{}
	return
}

// GetTeamDetails implements account.GetTeamDetails
func (s *Server) GetTeamDetails(ctx context.Context, in *GetTeamDetailsRequest) (out *GetTeamDetailsReply, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &GetTeamDetailsReply{}
	return
}

// DeleteTeam implements account.DeleteTeam
func (s *Server) DeleteTeam(ctx context.Context, in *DeleteTeamRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &pb.Empty{}
	return
}

// AddTeamMemberships implements account.AddTeamMemberships
func (s *Server) AddTeamMemberships(ctx context.Context, in *AddTeamMembershipsRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	out = &pb.Empty{}
	return
}

// DeleteTeamMemberships implements account.DeleteTeamMemberships
func (s *Server) DeleteTeamMemberships(ctx context.Context, in *DeleteTeamMembershipsRequest) (out *pb.Empty, err error) {
	if in.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if len(in.Members) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "members are mandatory")
	}
	out = &pb.Empty{}
	return
}

// GrantPermission implements account.GrantPermission
func (s *Server) GrantPermission(ctx context.Context, in *GrantPermissionRequest) (out *pb.Empty, err error) {
	//if in.Team == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	//}
	//if in.Organization == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	//}
	//if in.Level == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "permission level is mandatory")
	//}
	//if in.ResourceId == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	//}
	out = &pb.Empty{}
	return
}

// ListPermission implements account.ListPermission
func (s *Server) ListPermission(ctx context.Context, in *ListPermissionRequest) (out *ListPermissionReply, err error) {
	if in.Team != "" && in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &ListPermissionReply{}
	return
}

// EditPermission implements account.EditPermission
func (s *Server) EditPermission(ctx context.Context, in *EditPermissionRequest) (out *pb.Empty, err error) {
	//if in.Team == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	//}
	//if in.Organization == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	//}
	//if in.Level == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "permission level is mandatory")
	//}
	//if in.ResourceId == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	//}
	out = &pb.Empty{}
	return
}

// RevokePermission implements account.RevokePermission
func (s *Server) RevokePermission(ctx context.Context, in *RevokePermissionRequest) (out *pb.Empty, err error) {
	//if in.Team == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	//}
	//if in.Organization == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	//}
	//if in.ResourceId == "" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	//}
	out = &pb.Empty{}
	return
}

// TransferPermissionOwnership implements account.TransferPermissionOwnership
func (s *Server) TransferPermissionOwnership(ctx context.Context, in *TransferPermissionOwnershipRequest) (out *pb.Empty, err error) {
	if in.Team == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "team name is mandatory")
	}
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	if in.ResourceId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "resource id is mandatory")
	}
	out = &pb.Empty{}
	return
}
