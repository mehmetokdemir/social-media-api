package main

import (
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/auth"
	"github.com/mehmetokdemir/social-media-api/internal/app/comment"
	"github.com/mehmetokdemir/social-media-api/internal/app/friendship"
	"github.com/mehmetokdemir/social-media-api/internal/app/post"
	"github.com/mehmetokdemir/social-media-api/internal/app/user"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"github.com/mehmetokdemir/social-media-api/internal/logger"
	"github.com/mehmetokdemir/social-media-api/internal/postgres"
	"github.com/mehmetokdemir/social-media-api/server"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	appConfig := config.NewConfig()
	appConfig.Print()

	zapLogger, err := logger.NewZapLoggerForEnv(appConfig.AppEnv, 0)
	if err != nil {
		return err
	}

	db, err := postgres.New(appConfig)
	if err != nil {
		return err
	}

	userRepository := user.NewRepository(db, zapLogger)
	if err = userRepository.Migration(); err != nil {
		return nil
	}
	userService := user.NewUserService(userRepository, zapLogger, appConfig)
	userHandler := user.NewHttpHandler(userService, zapLogger, appConfig.JwtATPrivateKey)

	authService := auth.NewAuthService(userRepository, zapLogger, appConfig)
	authHandler := auth.NewHttpHandler(authService, zapLogger, appConfig.JwtATPrivateKey)

	friendshipRepository := friendship.NewRepository(db, zapLogger)
	if err = friendshipRepository.Migration(); err != nil {
		return nil
	}

	// TODO: IGNORE userRepository just check user with id and return bool
	friendshipService := friendship.NewFriendshipService(userRepository, friendshipRepository, zapLogger, appConfig)
	friendshipHandler := friendship.NewHttpHandler(friendshipService, zapLogger, appConfig.JwtATPrivateKey)

	postRepository := post.NewRepository(db, zapLogger)
	if err = postRepository.Migration(); err != nil {
		return nil
	}

	postService := post.NewPostService(postRepository, zapLogger, appConfig)
	postHandler := post.NewHttpHandler(postService, zapLogger, appConfig.JwtATPrivateKey)

	commentRepository := comment.NewRepository(db, zapLogger)
	if err = commentRepository.Migration(); err != nil {
		return nil
	}

	commentService := comment.NewCommentService(commentRepository, zapLogger, appConfig)
	commentHandler := comment.NewHttpHandler(commentService, zapLogger, appConfig.JwtATPrivateKey)

	appServer := server.New([]server.Handler{
		userHandler,
		authHandler,
		friendshipHandler,
		postHandler,
		commentHandler,
	}, appConfig, zapLogger)

	fmt.Println("server is start")
	return appServer.Start()
}
