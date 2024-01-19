package main

import (
	"crud/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//CRUD - Create, Read, Update, Delete
	//Create - método POST
	//Read - método GET
	//Update - método PUT
	//Delete - método DELETE

	router := mux.NewRouter()
	//Usa a rota para criar um usuario (método POST)
	router.HandleFunc("/usuarios", servidor.CriarUsuario).Methods(http.MethodPost)
	//Usa a rota para buscar usuarios (método GET)
	router.HandleFunc("/usuarios", servidor.BuscarUsuarios).Methods(http.MethodGet)
	// Usa a rota para buscar um usuario específico (método GET e passa o id como parâmetro)
	router.HandleFunc("/usuarios/{id}", servidor.BuscarUsuario).Methods(http.MethodGet)
	// Usa a rota para atualizar os dados de um usuario no banco de dados
	router.HandleFunc("/usuarios/{id}", servidor.AtualizarUsuario).Methods(http.MethodPut)
	// Usa a rota para deletar um usuario no banco de dados
	router.HandleFunc("/usuarios/{id}", servidor.DeletarUsuario).Methods(http.MethodDelete)

	fmt.Println("Servidor rodando na porta 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
