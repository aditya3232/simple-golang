package echo

import (
	"errors"
	"net/http"
	"simple-golang/internal/adapter/inbound/echo/request"
	"simple-golang/internal/adapter/inbound/echo/response"
	"simple-golang/internal/domain/entity"
	"simple-golang/internal/domain/service"
	"simple-golang/internal/port/inbound"
	"simple-golang/util"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type userHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) inbound.UserHandlerInterface {
	return &userHandler{userService: userService}
}

func (u *userHandler) SignIn(c echo.Context) error {
	var (
		req        = request.SignInRequest{}
		resp       = response.DefaultResponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-1] SignIn", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-2] SignIn", err)
	}

	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}
	user, token, err := u.userService.SignIn(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-3] SignIn", err)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-4] SignIn", err)
	}

	respSignIn.AccessToken = token
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.Phone = user.Phone

	resp.Message = "Success"
	resp.Data = respSignIn

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdatePassword(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		req  = request.UpdatePasswordRequest{}
		ctx  = c.Request().Context()
	)

	user, ok := c.Get("user").(entity.JwtUserData)
	if !ok {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdatePassword", err)
	}

	if user.UserID == 0 {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdatePassword", err)
	}

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-1] UpdatePassword", err)
	}

	if err := c.Validate(req); err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-2] UpdatePassword", err)
	}

	if req.NewPassword != req.ConfirmPassword {
		err := errors.New("new password and confirm password does not match")
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-3] UpdatePassword", err)
	}

	reqEntity := entity.UserEntity{
		ID:       user.UserID,
		Token:    user.Token,
		Password: req.NewPassword,
	}

	err := h.userService.UpdatePassword(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-4] UpdatePassword", errNotFound)
		}

		if err.Error() == "401" {
			errUnauthorized := errors.New("token expired or invalid")
			return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-5] UpdatePassword", errUnauthorized)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-6] UpdatePassword", err)
	}

	resp.Data = nil
	resp.Message = "Password updated successfully"

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)

	user, ok := c.Get("user").(entity.JwtUserData)
	if !ok {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] DeleteUser", err)
	}

	if user.UserID == 0 {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] DeleteUser", err)
	}

	idParamStr := c.Param("id")
	if idParamStr == "" {
		err := errors.New("missing or invalid user ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-1] DeleteUser", err)
	}

	id, err := util.StringToInt64(idParamStr)
	if err != nil {
		err := errors.New("invalid customer ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] DeleteUser", err)
	}

	err = h.userService.DeleteUser(ctx, id)
	if err != nil {
		log.Infof("[UserHandler-3] DeleteUser: %v", err)
		if err.Error() == "404" {
			errNotFOund := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-4] DeleteUser", errNotFOund)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-5] DeleteUser", err)
	}

	resp.Message = "Customer deleted successfully"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
		req  = request.UserUpdateRequest{}
	)

	user, ok := c.Get("user").(entity.JwtUserData)
	if !ok {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdateUser", err)
	}

	if user.UserID == 0 {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdateUser", err)
	}

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-1] UpdateUser", err)
	}

	if err := c.Validate(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] UpdateUser", err)
	}

	idParamStr := c.Param("id")
	if idParamStr == "" {
		err := errors.New("missing or invalid user ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-3] UpdateUser", err)
	}

	id, err := util.StringToInt64(idParamStr)
	if err != nil {
		err := errors.New("invalid user ID")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-4] UpdateUser", err)
	}

	reqEntity := entity.UserEntity{
		ID:      id,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	err = h.userService.UpdateUser(ctx, reqEntity)
	if err != nil {
		log.Errorf("[UserHandler-5] UpdateUser: %v", err)
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-6] UpdateUser", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-7] UpdateUser", err)
	}

	resp.Message = "Success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) CreateUserAccount(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
		req  = request.SignUpRequest{}
	)

	if err := c.Bind(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-1] CreateUserAccount", err)
	}

	if err := c.Validate(&req); err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] CreateUserAccount", err)
	}

	if req.Password != req.PasswordConfirmation {
		err := errors.New("password and confirm password does not match")
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UserHandler-3] CreateUserAccount", err)
	}

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
	}

	err := h.userService.CreateUserAccount(ctx, reqEntity)
	if err != nil {
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-4] CreateUserAccount", err)
	}

	resp.Message = "success"
	resp.Data = nil

	return c.JSON(http.StatusCreated, resp)
}

func (h *userHandler) GetUserByID(c echo.Context) error {
	var (
		resp     = response.DefaultResponse{}
		ctx      = c.Request().Context()
		respUser = response.UserResponse{}
	)

	idParam := c.Param("id")
	if idParam == "" {
		err := errors.New("id invalid")
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-1] GetUserByID", err)
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UserHandler-2] GetUserByID", err)
	}

	result, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		log.Errorf("[UserHandler-3] GetUserByID: %v", err)
		if err.Error() == "404" {
			errNotFound := errors.New("user not found")
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-4] GetUserByID", errNotFound)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-5] GetUserByID", err)
	}

	respUser.ID = result.ID
	respUser.Name = result.Name
	respUser.Email = result.Email
	respUser.Phone = result.Phone
	respUser.Address = result.Address

	resp.Data = respUser
	resp.Message = "success get user by id"

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) GetUser(c echo.Context) error {
	var (
		resp     = response.DefaultResponseWithPaginations{}
		ctx      = c.Request().Context()
		respUser = []response.UserResponse{}
	)

	user, ok := c.Get("user").(entity.JwtUserData)
	if !ok {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdateUser", err)
	}

	if user.UserID == 0 {
		err := errors.New("data token not valid")
		return response.RespondWithError(c, http.StatusUnauthorized, "[UserHandler-1] UpdateUser", err)
	}

	search := c.QueryParam("search")
	orderBy := "created_at"
	if c.QueryParam("order_by") != "" {
		orderBy = c.QueryParam("order_by")
	}

	orderType := c.QueryParam("order_type")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}

	pageStr := c.QueryParam("page")
	var page int64 = 1
	if pageStr != "" {
		page, _ = util.StringToInt64(pageStr)
		if page <= 0 {
			page = 1
		}
	}

	limitStr := c.QueryParam("limit")
	var limit int64 = 10
	if limitStr != "" {
		limit, _ = util.StringToInt64(limitStr)
		if limit <= 0 {
			limit = 10
		}
	}

	reqEntity := entity.QueryParamEntity{
		Search:    search,
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
	}

	results, countData, totalPages, err := h.userService.GetUser(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			return response.RespondWithError(c, http.StatusNotFound, "[UserHandler-1] Getuser", err)
		}
		return response.RespondWithError(c, http.StatusInternalServerError, "[UserHandler-2] Getuser", err)
	}

	for _, val := range results {
		respUser = append(respUser, response.UserResponse{
			ID:      val.ID,
			Name:    val.Name,
			Email:   val.Email,
			Phone:   val.Phone,
			Address: val.Address,
		})
	}

	resp.Message = "Data retrieved successfully"
	resp.Data = respUser
	resp.Pagination = &response.Pagination{
		Page:       page,
		TotalCount: countData,
		Limit:      limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, resp)
}
