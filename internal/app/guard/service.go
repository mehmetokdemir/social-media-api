package guard

type IGuardService interface {
	CheckTokenInBlacklist(token string) bool
}

type guardService struct {
	guardRepository IGuardRepository
}

func NewGuardService(guardRepository IGuardRepository) IGuardService {
	return &guardService{guardRepository: guardRepository}
}

func (s *guardService) CheckTokenInBlacklist(token string) bool {
	return s.guardRepository.CheckTokenInBlacklist(token)
}
