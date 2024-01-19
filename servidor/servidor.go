package servidor

import (
	"crud/conexao"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Função para criar usuarios (POST)

// CriarUsuario insere um usuário no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuário para struct!"))
		return
	}

	db, erro := conexao.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close()

	//prepare statement
	statement, erro := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	//executar o statement
	//USUARIO.nOME E USUARIO.EMAIL SÃO OS PARÂMETROS DO STATEMENT e precisam estar em ordem de acordo com o insert into usuarios (nome, email) values (?, ?)
	insercao, erro := statement.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte("Erro ao executar o statement!"))
		return
	}

	//Sempre que é inserido um novo usuário, é retornado o ID do usuário inserido como constume do MySQL
	//Para pegar o id inserido deve fazer:
	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao obter o ID inserido!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! ID: %d", idInserido)))
	//Se passar por todas as etapas acima, o usuário foi inserido com sucesso

	//o w.Writeheader(http.StatusCreated) é para retornar o status code 201 - created
	// Status code:
	// Deve fazer o status code em requisição http, ele diz se a requisicao foi bem sucedida ou não
	// 200 - sucesso
	// 400 - bad request
	// 500 - internal server error
	// 404 - not found
	// 401 - unauthorized
	// 403 - forbidden
	// 201 - created
	// 204 - no content
	// 304 - not modified
	// 301 - moved permanently
	// 302 - moved
	// 307 - temporary redirect
	// 308 - permanent redirect
	// 409 - conflict
	// 412 - precondition failed
	// 413 - payload too large
	//etc

}

//BUSCAR USUARIOS (GET)

// Função para buscar usuarios
// BuscarUsuarios traz todos os usuários salvos no banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := conexao.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close()

	//Select * from usuarios

	linhas, erro := db.Query("select * from usuarios")
	if erro != nil {
		w.Write([]byte("Erro ao buscar os usuários!"))
		return
	}
	defer linhas.Close()
	//sempre deve fechar as linhas

	var usuarios []usuario
	//linhas.Next() é para percorrer todas as linhas, para cada linha que o for linhas.Next percorrer ele vai executar uma iteração
	for linhas.Next() {
		var usuario usuario
		//linhas.Scan é para escanear os dados da linha e colocar dentro da variável usuario
		//usuario.ID, usuario.Nome, usuario.Email são os dados que estão no banco de dados que foram colocados pela gente, e deve ser colocado na ordem que está pois foi essa ordem em que criamos no banco de dados, primeiro id depois nome e dps email
		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear o usuário!"))
			return
		}
		usuarios = append(usuarios, usuario)
	}
	w.WriteHeader(http.StatusOK)

	//transformar o slice de usuarios em json
	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		//json.NewEncoder(w) é para codificar o slice de usuarios em json
		//Encode(usuarios) é como vai ser passado os dados
		w.Write([]byte("Erro ao converter os usuários para JSON!"))
		return
	}
}

// Função para buscar usuario
// BuscarUsuario tras um usuário específico salvo no banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {

	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parametro id para número inteiro"))
		return
	}

	db, erro := conexao.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}

	defer db.Close()

	linha, erro := db.Query("select * from usuarios where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao buscar o usuário!"))
		return
	}

	defer linha.Close()

	var usuario usuario
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear o usuário!"))
			return
		}
	}

	//função para ver o status se está ok
	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuário para JSON!"))
		return
	}
}

// Atualizar dados (PUT)

// AtualizarUsuario altera os dados do usuario no banco de dados
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {

	parametros := mux.Vars(r) //o r é o request

	//convertendo a string ID para número inteiro
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parametro id em número inteiro"))
		return
	}

	//até aqui já foi lido o parametro, já tem o id do usuario que ele está buscando, só não tem o corpo da requisição. Primeiro vai ser lido o corpo da requisição para depois abrir o banco

	//lendo o corpo da requisição
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição"))
		return
	}

	//puxando o slice de usuario com a var usuario
	var usuario usuario

	//aqui vai ser convertido o corpo da requisição e endereço de memória do usuario(&usuario) para struct
	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuário para struct"))
		return
	}

	//conectando com o banco de dados(abrindo ele)
	db, erro := conexao.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}

	defer db.Close()

	//criar o statement para proteger os dados do usuario
	statement, erro := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	//para atualizar o usuario
	//Ao invés de passar usuario.ID, deve passar apenas ID que é a variável atualizada que foi convertida em número inteiro
	if _, erro := statement.Exec(usuario.Nome, usuario.Email, ID); erro != nil {
		w.Write([]byte("Erro ao atualizar o usuário!"))
		return
	}

	//para aparecer o status de erro caso tenha
	w.WriteHeader(http.StatusNoContent)
}

// DELETAR usuario (DELETE)

//DeletarUsuario é uma função que remove um usuario do banco de dados

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {

	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parametro id em um número inteiro"))
		return
	}

	db, erro := conexao.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao remover usuário!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
