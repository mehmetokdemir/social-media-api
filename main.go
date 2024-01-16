package main

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"github.com/mehmetokdemir/social-media-api/internal/app/auth"
	"github.com/mehmetokdemir/social-media-api/internal/app/cdn"
	"github.com/mehmetokdemir/social-media-api/internal/app/comment"
	"github.com/mehmetokdemir/social-media-api/internal/app/friendship"
	"github.com/mehmetokdemir/social-media-api/internal/app/guard"
	"github.com/mehmetokdemir/social-media-api/internal/app/like"
	"github.com/mehmetokdemir/social-media-api/internal/app/post"
	"github.com/mehmetokdemir/social-media-api/internal/app/transaction"
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

	cloudinaryClient, err := cloudinary.NewFromURL(fmt.Sprintf("cloudinary://%s:%s@%s", appConfig.CloudinaryApiKey, appConfig.CloudinaryApiSecret, appConfig.CloudinaryCloudName))
	if err != nil {
		return err
	}
	cdnService := cdn.NewCdnService(cloudinaryClient)

	blackListRepository := auth.NewRepository(db, zapLogger)
	if err = blackListRepository.Migration(); err != nil {
		return nil
	}

	guardRepository := guard.NewRepository(db, zapLogger)
	guardService := guard.NewGuardService(guardRepository)

	userRepository := user.NewRepository(db, zapLogger)
	if err = userRepository.Migration(); err != nil {
		return nil
	}
	userService := user.NewUserService(userRepository, cdnService, zapLogger, appConfig)
	userHandler := user.NewHttpHandler(guardService, userService, zapLogger, appConfig.JwtATPrivateKey)

	authService := auth.NewAuthService(userService, blackListRepository, zapLogger, appConfig)
	authHandler := auth.NewHttpHandler(guardService, authService, zapLogger, appConfig.JwtATPrivateKey)

	friendshipRepository := friendship.NewRepository(db, zapLogger)
	if err = friendshipRepository.Migration(); err != nil {
		return nil
	}

	friendshipService := friendship.NewFriendshipService(friendshipRepository, zapLogger, appConfig)
	friendshipHandler := friendship.NewHttpHandler(guardService, friendshipService, zapLogger, appConfig.JwtATPrivateKey)

	commentRepository := comment.NewRepository(db, zapLogger)
	if err = commentRepository.Migration(); err != nil {
		return nil
	}
	commentService := comment.NewCommentService(commentRepository, cdnService, zapLogger, appConfig)
	commentHandler := comment.NewHttpHandler(guardService, commentService, zapLogger, appConfig.JwtATPrivateKey)

	likeRepository := like.NewRepository(db, zapLogger)
	if err = likeRepository.Migration(); err != nil {
		return nil
	}
	likeService := like.NewLikeService(likeRepository, zapLogger, appConfig)
	likeHandler := like.NewHttpHandler(guardService, likeService, zapLogger, appConfig.JwtATPrivateKey)

	transactionService := transaction.NewTransactionService(db)
	postRepository := post.NewRepository(db, zapLogger)
	if err = postRepository.Migration(); err != nil {
		return nil
	}
	postService := post.NewPostService(postRepository, cdnService, transactionService, commentService, likeService, zapLogger, appConfig)
	postHandler := post.NewHttpHandler(guardService, postService, zapLogger, appConfig.JwtATPrivateKey)

	appServer := server.New([]server.Handler{
		userHandler,
		authHandler,
		friendshipHandler,
		postHandler,
		commentHandler,
		likeHandler,
	}, appConfig, zapLogger)

	fmt.Println("server is start")
	return appServer.Start()
}
