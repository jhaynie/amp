package account

import (
	pb "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// Server is used to implement account.AccountServer
type Server struct{}

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (out *SignUpReply, err error) {
	return
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (out *pb.Empty, err error) {
	return
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (out *pb.Empty, err error) {
	return
}

// PasswordChange implements account.PasswordChange
func (s *Server) PasswordChange(ctx context.Context, in *PasswordChangeRequest) (out *pb.Empty, err error) {
	return
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (out *LogInReply, err error) {
	return
}

// List implements account.List
func (s *Server) List(ctx context.Context, in *ListRequest) (out *ListReply, err error) {
	return
}

// Switch implements account.Switch
func (s *Server) Switch(ctx context.Context, in *SwitchRequest) (out *pb.Empty, err error) {
	return
}

// GetDetails implements account.GetDetails
func (s *Server) GetDetails(ctx context.Context, in *GetDetailsRequest) (out *GetDetailsReply, err error) {
	return
}

// Edit implements account.Edit
func (s *Server) Edit(ctx context.Context, in *EditRequest) (out *EditRequestReply, err error) {
	return
}

// Delete implements account.Delete
func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (out *pb.Empty, err error) {
	return
}

// GetOrganizations implements account.GetOrganizations
func (s *Server) GetOrganizations(ctx context.Context, in *GetOrganizationsRequest) (out *GetOrganizationsReply, err error) {
	return
}

// GetTeams implements account.GetTeams
func (s *Server) GetTeams(ctx context.Context, in *GetTeamsRequest) (out *GetTeamsReply, err error) {
	return
}
