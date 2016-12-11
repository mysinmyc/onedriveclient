package auth

//Decouple authentication
type AuthenticationProvider interface {
	GetAuthenticationToken() (*AuthenticationToken, error)
	GetApplicationInfo() (ApplicationInfo, error)
}

//ApplicationInfo returns application informations required for authentication
type ApplicationInfo struct {
	ClientID     string
	ClientSecret string
	Scope        []string
	RedirectURI  string
}

//StaticAuthenticationInfo holder for static athentication data
type StaticAuthenticationInfo struct {
	AuthenticationToken *AuthenticationToken `json:"authenticationToken"`
	ApplicationInfo     *ApplicationInfo     `json:"applicationInfo"`
}

func (vSelf *StaticAuthenticationInfo) GetAuthenticationToken() (*AuthenticationToken, error) {
	return vSelf.AuthenticationToken, nil
}

func (vSelf *StaticAuthenticationInfo) GetApplicationInfo() (ApplicationInfo, error) {
	return *vSelf.ApplicationInfo, nil
}
