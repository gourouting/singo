package model

// Run database migrations.

func migration() {
	// Auto migration mode.
	_ = DB.AutoMigrate(&User{})
}
