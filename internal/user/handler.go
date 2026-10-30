package user

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	DB                *sql.DB
	PrivateKeyPath    string
	AdminBootstrapKey string
}

type credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var input credentials
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "invalid request body")
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if input.Name == "" || !strings.Contains(input.Email, "@") || len(input.Password) < 8 {
		return fiber.NewError(400, "name, valid email and password of at least 8 characters are required")
	}
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return fiber.NewError(500, "could not secure password")
	}
	var id string
	err = h.DB.QueryRowContext(c.Context(),
		`INSERT INTO users(name,email,password_hash,role) VALUES($1,$2,$3,'customer') RETURNING id`,
		input.Name, input.Email, passwordHash).Scan(&id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return fiber.NewError(409, "email already registered")
		}
		return fiber.NewError(500, "could not create account")
	}
	return c.Status(201).JSON(fiber.Map{"id": id, "name": input.Name, "email": input.Email})
}

func (h *Handler) BootstrapAdmin(c *fiber.Ctx) error {
	if h.AdminBootstrapKey == "" || c.Get("X-Admin-Bootstrap-Key") != h.AdminBootstrapKey {
		return fiber.NewError(401, "valid admin bootstrap key required")
	}
	var input credentials
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "invalid request body")
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if input.Name == "" || !strings.Contains(input.Email, "@") || len(input.Password) < 12 {
		return fiber.NewError(400, "name, valid email and password of at least 12 characters are required")
	}
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return fiber.NewError(500, "could not secure password")
	}
	tx, err := h.DB.BeginTx(c.Context(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fiber.NewError(500, "could not initialize administrator")
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(c.Context(), `SELECT pg_advisory_xact_lock(20261030)`); err != nil {
		return fiber.NewError(500, "could not initialize administrator")
	}
	var count int
	if err = tx.QueryRowContext(c.Context(), `SELECT count(*) FROM users WHERE role='admin' AND deleted_at IS NULL`).Scan(&count); err != nil {
		return fiber.NewError(500, "could not initialize administrator")
	}
	if count > 0 {
		return fiber.NewError(409, "administrator is already initialized")
	}
	var id string
	if err = tx.QueryRowContext(c.Context(), `INSERT INTO users(name,email,password_hash,role) VALUES($1,$2,$3,'admin') RETURNING id`, input.Name, input.Email, passwordHash).Scan(&id); err != nil {
		return fiber.NewError(409, "administrator account could not be created")
	}
	if err = tx.Commit(); err != nil {
		return fiber.NewError(500, "could not finalize administrator")
	}
	return c.Status(201).JSON(fiber.Map{"id": id, "name": input.Name, "email": input.Email, "role": "admin"})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var input credentials
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "invalid request body")
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	var id, name, passwordHash, role string
	err := h.DB.QueryRowContext(c.Context(),
		`SELECT id,name,password_hash,role FROM users WHERE email=$1 AND deleted_at IS NULL`,
		input.Email).Scan(&id, &name, &passwordHash, &role)
	if err != nil || !checkPassword(passwordHash, input.Password) {
		return fiber.NewError(401, "invalid email or password")
	}
	key, err := loadPrivateKey(h.PrivateKeyPath)
	if err != nil {
		return fiber.NewError(500, "token signing is unavailable")
	}
	token, err := issueToken(key, id, input.Email, role)
	if err != nil {
		return fiber.NewError(500, "could not create access token")
	}
	return c.JSON(fiber.Map{"access_token": token, "token_type": "Bearer", "expires_in": 86400,
		"user": fiber.Map{"id": id, "name": name, "email": input.Email, "role": role}})
}

func (h *Handler) Me(c *fiber.Ctx) error {
	id := c.Locals("userID").(string)
	var name, email, role string
	err := h.DB.QueryRowContext(c.Context(), `SELECT name,email,role FROM users WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&name, &email, &role)
	if err != nil {
		return fiber.NewError(404, "user not found")
	}
	return c.JSON(fiber.Map{"id": id, "name": name, "email": email, "role": role})
}
