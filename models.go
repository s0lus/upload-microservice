package main

type Image struct {
	ID 		int  	`gorm:"not null;type:serial;primary_key;unique"`
	Name 	string	`gorm:"not null;type:varchar(60);unique"`
	Width 	int		`gorm:"not null;type:smallserial"`
	Height 	int		`gorm:"not null;type:smallserial"`
}

type Config struct {
	DataBase struct {
		Dialect 		   string   `json:"Dialect"`
		ConnectionData 	   string	`json:"ConnectionData"`
		IdleConnections    int	    `json:"IdleConnections"`
		MaxOpenConnections int      `json:"MaxOpenConnections"`
	}
}
