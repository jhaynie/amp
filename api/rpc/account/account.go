package account

import (
	"github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/data/storage"
	"github.com/dgrijalva/jwt-go"
	pb "github.com/golang/protobuf/ptypes/empty"
	"github.com/hlandau/passlib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"os"
	"time"
)

// TODO: this MUST NOT be public
// TODO: find a way to store this key secretly
var secretKey = []byte("&kv@l3go-f=@^*@ush0(o5*5utxe6932j9di+ume=$mkj%d&&9*%k53(bmpksf&!c2&zpw$z=8ndi6ib)&nxms0ia7rf*sj9g8r4")

type accountClaims struct {
	AccountID string `json:"AccountID"`
	jwt.StandardClaims
}

// Server is used to implement account.AccountServer
type Server struct {
	accounts account.Interface
}

// NewServer instantiates account.Server
func NewServer(store storage.Interface) *Server {
	return &Server{accounts: account.NewStore(store)}
}

// SignUp implements account.SignUp
func (s *Server) SignUp(ctx context.Context, in *SignUpRequest) (*SignUpReply, error) {
	// Validate input
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Check if account already exists
	alreadyExists, err := s.accounts.GetAccountByUserName(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if alreadyExists != nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "account already exists")
	}

	// Hash password
	passwordHash, err := passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Create the new account
	account := &schema.Account{
		Email:        in.Email,
		UserName:     in.UserName,
		Type:         in.AccountType,
		IsVerified:   false,
		PasswordHash: passwordHash,
	}
	id, err := s.accounts.CreateAccount(ctx, account)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "storage error")
	}
	log.Println("Successfully created account", in.UserName)

	// Forge the verification token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accountClaims{
		id, // The token contains the account id to verify
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    os.Args[0],
		},
	})

	// Sign the token
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// TODO: send confirmation email with token

	return &SignUpReply{Token: ss}, nil
}

// Verify implements account.Verify
func (s *Server) Verify(ctx context.Context, in *VerificationRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Validate the token
	token, err := jwt.ParseWithClaims(in.Token, &accountClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if !token.Valid {
		return &pb.Empty{}, grpc.Errorf(codes.InvalidArgument, "invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*accountClaims)
	if !ok {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, "invalid claims")
	}

	// Activate the account
	account, err := s.accounts.GetAccount(ctx, claims.AccountID)
	if err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	account.IsVerified = true
	if err := s.accounts.UpdateAccount(ctx, account); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully verified account", account.UserName)

	return &pb.Empty{}, nil
}

// Login implements account.Login
func (s *Server) Login(ctx context.Context, in *LogInRequest) (*LogInReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the account
	account, err := s.accounts.GetAccountByUserName(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if account == nil {
		return nil, grpc.Errorf(codes.NotFound, "account not found")
	}
	if !account.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "account not verified")
	}

	// Check password
	_, err = passlib.Verify(in.Password, account.PasswordHash)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	// Forge the authentication token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accountClaims{
		account.Id, // The token contains the account id
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    os.Args[0],
		},
	})

	// Sign the token
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully login for account", account.UserName)

	return &LogInReply{Token: ss}, nil
}

// PasswordReset implements account.PasswordReset
func (s *Server) PasswordReset(ctx context.Context, in *PasswordResetRequest) (*PasswordResetReply, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Get the account
	account, err := s.accounts.GetAccountByUserName(ctx, in.UserName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if account == nil {
		return nil, grpc.Errorf(codes.NotFound, "account not found")
	}
	// TODO: Do we need the account to be verified?
	if !account.IsVerified {
		return nil, grpc.Errorf(codes.FailedPrecondition, "account not verified")
	}

	// Forge the password reset token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accountClaims{
		account.Id, // The token contains the account id to reset
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    os.Args[0],
		},
	})

	// Sign the token
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully reset password for account", account.UserName)

	// TODO: send password reset email with token

	return &PasswordResetReply{Token: ss}, nil
}

// PasswordSet implements account.PasswordSet
func (s *Server) PasswordSet(ctx context.Context, in *PasswordSetRequest) (*pb.Empty, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Validate the token
	token, err := jwt.ParseWithClaims(in.Token, &accountClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if !token.Valid {
		return &pb.Empty{}, grpc.Errorf(codes.InvalidArgument, "invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*accountClaims)
	if !ok {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, "invalid claims")
	}

	// Get the account
	account, err := s.accounts.GetAccount(ctx, claims.AccountID)
	if err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}

	// Sets the new password
	passwordHash, err := passlib.Hash(in.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	account.PasswordHash = passwordHash
	if err := s.accounts.UpdateAccount(ctx, account); err != nil {
		return &pb.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	log.Println("Successfully set new password for account", account.UserName)

	return &pb.Empty{}, nil
}

// PasswordChange implements account.PasswordChange
func (s *Server) PasswordChange(ctx context.Context, in *PasswordChangeRequest) (out *pb.Empty, err error) {
	// TODO: check if account is verified
	return
}

// List implements account.List
func (s *Server) List(ctx context.Context, in *ListAccountRequest) (out *ListAccountReply, err error) {
	// TODO: check if account is verified
	//if in.Type != "individual" && in.Type != "organization" {
	//	return nil, grpc.Errorf(codes.InvalidArgument, "account type is mandatory")
	//}
	out = &ListAccountReply{}
	return
}

// Switch implements account.Switch
func (s *Server) Switch(ctx context.Context, in *SwitchRequest) (out *pb.Empty, err error) {
	// TODO: check if account is verified
	if in.UserName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &pb.Empty{}
	return
}

// GetDetails implements account.GetDetails
func (s *Server) GetDetails(ctx context.Context, in *GetAccountDetailsRequest) (out *GetAccountDetailsReply, err error) {
	// TODO: check if account is verified
	if in.UserName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	out = &GetAccountDetailsReply{}
	return
}

// Edit implements account.Edit
func (s *Server) Edit(ctx context.Context, in *EditAccountRequest) (out *pb.Empty, err error) {
	// TODO: check if account is verified
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
	// TODO: check if account is verified
	if in.UserName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "name is mandatory")
	}
	out = &pb.Empty{}
	return
}

// GetTeams implements account.GetTeams
func (s *Server) GetTeams(ctx context.Context, in *GetTeamsRequest) (out *GetTeamsReply, err error) {
	// TODO: check if account is verified
	out = &GetTeamsReply{}
	return
}

// AddOrganizationMemberships implements account.AddOrganizationMemberships
func (s *Server) AddOrganizationMemberships(ctx context.Context, in *AddOrganizationMembershipsRequest) (out *pb.Empty, err error) {
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
	if in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &ListTeamReply{}
	return
}

// EditTeam implements account.EditTeam
func (s *Server) EditTeam(ctx context.Context, in *EditTeamRequest) (out *pb.Empty, err error) {
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
	if in.Team != "" && in.Organization == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization name is mandatory")
	}
	out = &ListPermissionReply{}
	return
}

// EditPermission implements account.EditPermission
func (s *Server) EditPermission(ctx context.Context, in *EditPermissionRequest) (out *pb.Empty, err error) {
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
	// TODO: check if account is verified
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
