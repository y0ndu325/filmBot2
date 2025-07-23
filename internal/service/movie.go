package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"film_botV2/internal/database"
	"film_botV2/internal/models"
	"film_botV2/internal/tmdb"
)

var (
	ErrNoMovies     = errors.New("no movies in database")
	ErrInvalidIndex = errors.New("invalid movie index")
)


type MovieService struct {
	db *database.Database
	tmdb *tmdb.Client
}

func New(db *database.Database, tmdb *tmdb.Client) *MovieService {
	return &MovieService{db: db, tmdb: tmdb}
}

func (s *MovieService) AddMovie(title string) error {
	movie := &models.Movie{
		Title: title,
	}
	return s.db.AddMovie(movie)
}

func (s *MovieService) GetMovieDetailsWithYear(title, year string) (*models.Movie, error) {
	movie, err := s.db.GetMovieByTitle(title)
	if err == nil && movie.Overview != "" {
		return movie, nil
	}

	searchResults, err := s.tmdb.SearchMovie(title, year)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска фильма: %w", err)
	}

	if len(searchResults.Results) == 0 {
		return nil, fmt.Errorf("фильм не найден")
	}

	movieDetails, err := s.tmdb.GetMovieDetails(int(searchResults.Results[0].ID))
	if err != nil {
		return nil, fmt.Errorf("ошибка получения деталей фильма: %w", err)
	}

	movie = &models.Movie{
		Title:       movieDetails.Title,
		Overview:    movieDetails.Overview,
		PosterPath:  movieDetails.PosterPath,
		ReleaseDate: movieDetails.ReleaseDate,
	}

	if err := s.db.AddMovie(movie); err != nil {
		return nil, fmt.Errorf("ошибка добавления фильма в базу: %w", err)
	}

	return movie, nil
}

func (s *MovieService) AddMovieWithYear(title, year string) error{
	_, err := s.GetMovieDetailsWithYear(title, year)
	return err
}

func (s *MovieService) GetRandomMovie() (*models.Movie, error) {
	movies, err := s.db.GetMovies()
	if err != nil {
		return  nil, err
	}
	if len(movies) == 0 {
		return nil, ErrNoMovies
	}

	rand.Seed(time.Now().UnixNano())
    randomIndex := rand.Intn(len(movies))
    movie := &movies[randomIndex]

	if movie.Overview == "" || movie.PosterPath == "" {
		details, err := s.GetMovieDetailsWithYear(movie.Title, movie.ReleaseDate)
		if err == nil {
			movie.Overview = details.Overview
			movie.PosterPath = details.PosterPath
		}
	}

	if err := s.db.DeleteMovie(movie); err != nil{
		return nil, err
	}

	return movie, nil
}

func (s *MovieService) GetMovies() ([]models.Movie, error) {
	return s.db.GetMovies()
}

func (s *MovieService) DeleteMovie(index int) (*models.Movie, error) {
	movies, err := s.db.GetMovies()
	if err != nil {
		return  nil, err
	}
	if index < 1 || index > len(movies) {
		return nil, ErrInvalidIndex
	}

	movie := movies[index-1]
	if err := s.db.DeleteMovie(&movie); err != nil{
		return nil, err
	}
	return &movie, nil
}

func (s *MovieService) TMDB() *tmdb.Client {
	return s.tmdb
}