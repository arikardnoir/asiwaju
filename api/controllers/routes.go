package controllers

import "github.com/arikardnoir/asiwaju/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Shops routes
	s.Router.HandleFunc("/shops", middlewares.SetMiddlewareJSON(s.CreateShop)).Methods("POST")
	s.Router.HandleFunc("/shops", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetShops))).Methods("GET")
	s.Router.HandleFunc("/shops/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetShop))).Methods("GET")
	s.Router.HandleFunc("/shops/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateShop))).Methods("PUT")
	s.Router.HandleFunc("/shops/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteShop)).Methods("DELETE")

	//Addresses routes
	s.Router.HandleFunc("/address", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateAddress))).Methods("POST")
	s.Router.HandleFunc("/address", middlewares.SetMiddlewareJSON(s.GetAddresses)).Methods("GET")
	s.Router.HandleFunc("/address/{id}", middlewares.SetMiddlewareJSON(s.GetAddress)).Methods("GET")
	s.Router.HandleFunc("/address/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateAddress))).Methods("PUT")
	s.Router.HandleFunc("/address/{id}", middlewares.SetMiddlewareAuthentication(middlewares.SetMiddlewareAuthentication(s.DeleteAddress))).Methods("DELETE")

	//Collection Point route
	s.Router.HandleFunc("/collection-points", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCollectionPoint))).Methods("GET")
	s.Router.HandleFunc("/collection-points/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCollectionPoint))).Methods("GET")

}
