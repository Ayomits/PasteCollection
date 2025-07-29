package services

type AuthService interface {
	Register() error
}

type authService struct {
	configService ConfigService
}

func NewAuthService(configService ConfigService) AuthService {
	return &authService{configService: configService}
}

func (s *authService) Register() error {
	dbURL, err := s.configService.Get("DB_URL")
	if err != nil {
		return err
	}

	_ = dbURL

	return nil
}
