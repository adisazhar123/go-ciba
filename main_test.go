package go_ciba

func main() {
	// pollInterval := 5
	// ds := repository.NewInMemoryDataStore()

	// clientAppRepo repository.ClientApplicationRepositoryInterface,
	// 	userAccountRepo repository.UserAccountRepositoryInterface,
	// 	cibaSessionRepo repository.CibaSessionRepositoryInterface,
	// 	keyRepo repository.KeyRepositoryInterface,
	// 	notificationClient transport.NotificationInterface,
	// 	cibaGrant *grant.CibaGrant,
	// 	validateClientNotificationToken func(token string) bool,

	// cibaService := service.NewCibaService(
	// 	ds.GetClientApplicationRepository(),
	// 	ds.GetUserAccountRepository(),
	// 	ds.GetCibaSessionRepository(),
	// 	ds.GetKeyRepository(),
	// 	transport.NewFirebaseCloudMessaging("123123"),
	// 	grant.NewCustomCibaGrant(&pollInterval, &grant.GrantConfig{
	// 		Issuer:              "https://adisazhar.com",
	// 		IdTokenLifetime:     3600,
	// 		AccessTokenLifetime: 3600,
	// 	}),
	// 	func(token string) bool {
	// 		return token != ""
	// 	},
	// )
	//
	//
	// as := NewAuthorizationServer(ds)
	// as.AddService(cibaService)
	//
	// as.HandleCibaRequest()
}
