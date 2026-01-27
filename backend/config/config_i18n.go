// config_i18n.go - 配置文件国际化注释
// config_i18n.go - Configuration file internationalization comments

package config

// InitConfig 初始化配置 / Initialize configuration
func InitConfig() {
	// 加载 .env 文件 / Load .env file
	_ = godotenv.Load()

	// 初始化日志 / Initialize logging
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)

	// 初始化环境配置 / Initialize environment configuration
	InitEnvironment()
}

// InitDB 初始化数据库 / Initialize database
func InitDB() {
	var err error
	var dialector gorm.Dialector

	dbType := GetEnv("DB_TYPE", "sqlite")
	dbURL := GetEnv("DATABASE_URL", "")

	if dbType == "postgres" && dbURL != "" {
		// 使用 PostgreSQL / Use PostgreSQL
		dialector = postgres.Open(dbURL)
	} else {
		// 默认使用 SQLite（开发环境）/ Default to SQLite (development environment)
		dialector = sqlite.Open("sol-green.db")
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err)) // "Database connection failed"
	}

	// 自动迁移 / Auto migrate
	DB.AutoMigrate(
		&model.GreenBehavior{},
		&model.User{},
		&model.RewardRecord{},
		&model.Challenge{},
		&model.ChallengeParticipant{},
		&model.MarketingActivity{},
		&model.MarketingParticipant{},
		&model.SignInRecord{},
		&model.InviteRecord{},
	)

	Log.Info("数据库初始化成功") // "Database initialization successful"
}

// GetEnv 获取环境变量 / Get environment variable
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Now 返回当前时间 / Return current time
func Now() time.Time {
	return time.Now()
}
