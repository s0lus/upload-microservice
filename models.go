package main

type Image struct {
	ID 		int
	Name 	string
	Path 	string
	Width 	int
	Height 	int
}

type Config struct {
	DataBase struct {
		Dialect 		   string
		ConnectionData 	   string
		IdleConnections    int
		MaxOpenConnections int
	}
}
