package users

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"testapi/internal/entity"
	"testapi/json"
)

func Get(repository Repository, tokenAuth *jwtauth.JWTAuth) *chi.Mux {
	res := resource{repository: repository, tokenAuth: tokenAuth}
	r := chi.NewRouter()

	jwtMiddleware := []func(http.Handler) http.Handler{jwtauth.Verifier(tokenAuth), jwtauth.Authenticator}

	r.Post("/", res.createUser)

	r.With(jwtMiddleware...).Get("/", res.getUsers)
	r.With(jwtMiddleware...).Get("/{id}", res.getUserByID)
	r.With(jwtMiddleware...).Put("/{id}", res.updateUser)
	r.With(jwtMiddleware...).Delete("/{id}", res.deleteUser)
	return r
}

type resource struct {
	repository Repository
	tokenAuth  *jwtauth.JWTAuth
}

func (res resource) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := res.repository.Query(ctx, 0, 0)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Unable to get users")
		render.NoContent(w, r)
		return
	}

	if users == nil {
		log.WithFields(log.Fields{
			"err": "no users available",
		}).Info("Unable to get users")

		render.NoContent(w, r)
		return
	}

	for i := range users {
		users[i].Password = ""
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
}

func (res resource) getUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to get user by ID")

		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, json.H{"errors": fmt.Sprintf("user ID '%s' is invalid format", paramId)})
		return
	}

	user, err := res.repository.Get(ctx, int64(id))

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to get user by ID")

		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, json.H{"errors": fmt.Sprintf("user with ID '%s' not found", paramId)})
		return
	}

	user.Password = ""

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)

}

func (res resource) createUser(w http.ResponseWriter, r *http.Request) {
	var user = entity.User{}

	ctx := r.Context()

	if err := json.Decode(r, &user); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Unable to create user")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, json.H{"errors": "invalid json"})
		return
	}

	if err := user.ValidateCreateRequest(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Unable to create user")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, json.H{"errors": err})
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Unable to create user")

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, json.H{"errors": "unable to create user"})
		return
	}

	user.Password = string(pass)

	newUser, err := res.repository.Create(ctx, user)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Unable to create user")

		render.Status(r, http.StatusConflict)
		render.JSON(w, r, json.H{"errors": err.Error()})
		return
	}

	_, tokenString, err := res.tokenAuth.Encode(map[string]interface{}{"user_id": newUser.ID})

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("Unable to create user")

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, json.H{"errors": "unable to create user"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, json.H{"token": tokenString})

}

func (res resource) updateUser(w http.ResponseWriter, r *http.Request) {
	var user = entity.User{}

	paramId := chi.URLParam(r, "id")

	ctx := r.Context()

	id, err := strconv.Atoi(paramId)

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to update user")
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, json.H{"errors": fmt.Sprintf("User ID %s is invalid format", paramId)})
		return
	}

	_, err = res.repository.Get(ctx, int64(id))

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to update user")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, json.H{"errors": fmt.Sprintf("user with ID %s not found", paramId)})
		return
	}

	if err := json.Decode(r, &user); err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to update user")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, json.H{"errors": "invalid json"})
		return
	}

	if err := user.ValidateUpdateRequest(); err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to update user")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, json.H{"errors": err})
		return
	}

	updatedUser, err := res.repository.Update(ctx, int64(id), user)

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to update user")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, json.H{"errors": "unable to update user"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, updatedUser)

}

func (res resource) deleteUser(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")

	ctx := r.Context()

	id, err := strconv.Atoi(paramId)

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to delete user")

		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, json.H{"errors": fmt.Sprintf("user ID '%s' is invalid format", paramId)})
		return
	}

	_, err = res.repository.Get(ctx, int64(id))

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to delete user")

		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, json.H{"errors": fmt.Sprintf("user with ID '%s' not found", paramId)})
	}

	err = res.repository.Delete(ctx, int64(id))

	if err != nil {
		log.WithFields(log.Fields{
			"id":  paramId,
			"err": err,
		}).Info("Unable to delete user")

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, json.H{"errors": fmt.Sprintf("unable to delete user with ID '%s'", paramId)})
		return
	}

	w.WriteHeader(http.StatusOK)
}
