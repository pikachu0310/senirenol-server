package handler

import (
	"net/http"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RegisterUser godoc
// @Summary ユーザー登録
// @Description 初回起動時に匿名ユーザーを生成し、user_idを返します
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} RegisterUserResponse "登録されたユーザーのID"
// @Router /users [post]
func (h *Handler) RegisterUser(c echo.Context) error {
	id, err := h.repo.CreateUser(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return c.JSON(http.StatusOK, RegisterUserResponse{ID: id.String()})
}

type updateUserNameRequest struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

// UpdateUserName godoc
// @Summary ユーザー名更新
// @Description user_idに対してuser_nameを更新します
// @Tags users
// @Accept json
// @Produce json
// @Param body body updateUserNameRequest true "更新内容"
// @Success 200 {object} UpdateUserNameResponse "更新結果"
// @Failure 400 {object} ErrorResponse
// @Router /users/update [post]
func (h *Handler) UpdateUserName(c echo.Context) error {
	var req updateUserNameRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}
	if err := vd.ValidateStruct(&req,
		vd.Field(&req.UserID, vd.Required),
		vd.Field(&req.UserName, vd.Required, vd.RuneLength(1, 255)),
	); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if _, err := uuid.Parse(req.UserID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user_id").SetInternal(err)
	}
	if err := h.repo.UpdateUserName(c.Request().Context(), req.UserID, req.UserName); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return c.JSON(http.StatusOK, UpdateUserNameResponse{Status: "ok"})
}

// GetUser godoc
// @Summary ユーザー情報取得
// @Description 指定したIDのユーザー情報を取得します（名前のみ）
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID" format(uuid)
// @Success 200 {object} GetUserResponse "ユーザー情報"
// @Failure 400 {object} ErrorResponse
// @Router /users/{userID} [get]
func (h *Handler) GetUser(c echo.Context) error {
	userID := c.Param("userID")
	if _, err := uuid.Parse(userID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid userID").SetInternal(err)
	}
	u, err := h.repo.GetUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()).SetInternal(err)
	}
	return c.JSON(http.StatusOK, GetUserResponse{ID: u.ID, Name: u.Name})
}

// GetUserStats godoc
// @Summary ユーザー統計
// @Description プレイ回数など統計情報を返します
// @Tags users
// @Produce json
// @Param userID path string true "User ID" format(uuid)
// @Success 200 {object} UserStatsResponse "統計情報"
// @Failure 400 {object} ErrorResponse
// @Router /users/{userID}/stats [get]
func (h *Handler) GetUserStats(c echo.Context) error {
	uid := c.Param("userID")
	if _, err := uuid.Parse(uid); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid userID").SetInternal(err)
	}
	s, err := h.repo.GetUserStats(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return c.JSON(http.StatusOK, s)
}
