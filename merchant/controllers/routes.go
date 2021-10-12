package controllers

func (server *Server) initializeRoutes() {

	route := server.Router.HandleFunc

	route("/login", server.LoginController)

	route("/users", server.UserController)

	route("/products", server.ProductController)

}
