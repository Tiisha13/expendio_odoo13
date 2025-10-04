package routes

import (
	"expensio-backend/internal/config"
	"expensio-backend/internal/handler"
	"expensio-backend/internal/middleware"
	"expensio-backend/internal/repository"
	"expensio-backend/internal/service"
	"expensio-backend/pkg/ocr"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository()
	companyRepo := repository.NewCompanyRepository()
	expenseRepo := repository.NewExpenseRepository()
	approvalRepo := repository.NewApprovalRepository()
	approvalRuleRepo := repository.NewApprovalRuleRepository()
	ocrResultRepo := repository.NewOCRResultRepository()

	// Initialize services
	authService := service.NewAuthService(userRepo, companyRepo, cfg)
	userService := service.NewUserService(userRepo, companyRepo, cfg)
	expenseService := service.NewExpenseService(expenseRepo, userRepo, companyRepo, cfg)
	approvalService := service.NewApprovalService(approvalRepo, approvalRuleRepo, expenseRepo, userRepo, cfg)
	ocrService := ocr.NewOCRService(cfg)

	// Set approval service in expense service (to avoid circular dependency)
	expenseService.SetApprovalService(approvalService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService, cfg)
	expenseHandler := handler.NewExpenseHandler(expenseService, cfg)
	approvalHandler := handler.NewApprovalHandler(approvalService, cfg)
	ocrHandler := handler.NewOCRHandler(ocrService, ocrResultRepo, expenseService, cfg)

	// API v1 group
	api := app.Group("/api/v1")

	// Public routes (no authentication required)
	auth := api.Group("/auth")
	{
		auth.Post("/signup", authHandler.Signup)
		auth.Post("/login", authHandler.Login)
		auth.Post("/refresh", authHandler.RefreshToken)
	}

	// Protected routes (authentication required)
	protected := api.Group("", middleware.AuthMiddleware(cfg))
	{
		// Auth routes
		protected.Post("/auth/logout", authHandler.Logout)

		// User routes
		users := protected.Group("/users")
		{
			// Admin only routes
			users.Post("/", middleware.RoleMiddleware("admin"), userHandler.CreateUser)
			users.Put("/:id/role", middleware.RoleMiddleware("admin"), userHandler.UpdateUserRole)
			users.Delete("/:id", middleware.RoleMiddleware("admin"), userHandler.DeleteUser)

			// Admin and Manager routes
			users.Get("/", middleware.RoleMiddleware("admin", "manager"), userHandler.GetUsers)
			users.Put("/:id/manager", middleware.RoleMiddleware("admin", "manager"), userHandler.AssignManager)

			// All authenticated users
			users.Get("/:id", userHandler.GetUser)
		}

		// Expense routes
		expenses := protected.Group("/expenses")
		{
			// All authenticated users
			expenses.Post("/", expenseHandler.CreateExpense)
			expenses.Get("/", expenseHandler.GetExpenses)

			// Manager and Admin only - Must be BEFORE /:id route!
			expenses.Get("/pending", middleware.RoleMiddleware("admin", "manager"), expenseHandler.GetPendingExpenses)

			// All authenticated users - Dynamic routes should be last
			expenses.Get("/:id", expenseHandler.GetExpense)
			expenses.Put("/:id", expenseHandler.UpdateExpense)
			expenses.Delete("/:id", expenseHandler.DeleteExpense)
		}

		// Approval routes
		approvals := protected.Group("/approvals")
		{
			// Manager and Admin only
			approvals.Get("/pending", middleware.RoleMiddleware("admin", "manager"), approvalHandler.GetPendingApprovals)
			approvals.Post("/:id/approve", middleware.RoleMiddleware("admin", "manager"), approvalHandler.ApproveExpense)
			approvals.Post("/:id/reject", middleware.RoleMiddleware("admin", "manager"), approvalHandler.RejectExpense)

			// All authenticated users can view approval history
			approvals.Get("/history/:expenseId", approvalHandler.GetApprovalHistory)
		}

		// OCR routes
		ocr := protected.Group("/ocr")
		{
			// All authenticated users
			ocr.Post("/upload", ocrHandler.UploadReceipt)
		}
	}

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
		})
	})
}
