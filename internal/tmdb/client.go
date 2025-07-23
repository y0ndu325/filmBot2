package tmdb

import (
	"fmt"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type Client struct {
	apiKey string
	client *tmdb.Client
}

func (c *Client) GetMovieDetailsWithYear(title string, year string) (any, any) {
	panic("unimplemented")
}

func NewClient(apiKey string) (*Client, error) {
	c, err := tmdb.Init(apiKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации tmdb клиента: %w", err)
	}

	return &Client{
		apiKey: apiKey,
		client: c,
	}, nil
}

func (c *Client) SearchMovie(title, year string) (*tmdb.SearchMovies, error) {
	options := map[string]string{
		"query": title,
		"year":  year,
	}

	if len(year) != 4 {
		return nil, fmt.Errorf("год обязателен и должен иметь формат YYYY")
	}

	results, err := c.client.GetSearchMovies(title, options)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска фильма: %w", err)
	}
	return results, nil
}

func (c *Client) GetMovieDetails(movieId int) (*tmdb.MovieDetails, error) {
	options := map[string]string{}
	movie, err := c.client.GetMovieDetails(movieId, options)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения деталей фильма: %w", err)
	}
	return movie, nil
}

func (c *Client) GetImageURL(PosterPath string) string {
	if PosterPath == "" {
		return ""
	}
	baseURL := "https://image.tmdb.org/t/p/"
	size := "original"
	return fmt.Sprintf("%s%s%s", baseURL, size, PosterPath)
}
