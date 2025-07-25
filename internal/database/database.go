package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"film_botV2/internal/models"
)

type Database struct {
	db *gorm.DB
	tableExists bool
}

func NewDatabase(dsn string) (*Database, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	d := &Database{
		db: db,
	}

	if !db.Migrator().HasTable("movies") {
		if err := db.AutoMigrate(&models.Movie{}); err != nil {
			return nil, err
		}

		d.tableExists = true
	}

	return d, nil
}

func (d *Database) AddMovie(movie *models.Movie) error {
	return d.db.Create(movie).Error
}

func (d *Database) GetMovies() ([]models.Movie, error) {
	var movies []models.Movie
	return movies, d.db.Find(&movies).Error
}

func (d *Database) DeleteMovie(movie *models.Movie) error {
	return d.db.Delete(movie).Error
}

func (d *Database) GetMovieByTitle(title string) (*models.Movie, error) {
	var movie models.Movie
	err := d.db.Where("title = ?", title).First(&movie).Error
	if err != nil {
		return nil, err
	}

	return &movie, nil 
}

func (d *Database) GetRandomMovie() (*models.Movie, error) {
	var movie models.Movie
	err := d.db.Order("RANDOM()").First(&movie).Error
	if err != nil {
		return nil, err
	}

	return &movie, nil
}