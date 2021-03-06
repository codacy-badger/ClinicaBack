package pacientes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ClinicaBack/config"
	"ClinicaBack/model/paciente"
)

//DB ...
var DB = db.Con
var mensagemErro string

var pacientes []paciente.Paciente

// Adicionar ...
func Adicionar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Chamando rota adicionar paciente ...")
	w.Header().Set("Content-Type", "application/json")

	novoPaciente := paciente.Paciente{}

	err := json.NewDecoder(r.Body).Decode(&novoPaciente)
	mensagemErro = "erro_corpo"
	CheckErro(w, r, mensagemErro, err)

	query := "INSERT INTO paciente (nome, email, senha, data_nascimento, hospital, carteira, role, ativo) VALUES(?,?,?,?,?,?,?,?)"
	stmt, err := DB.Prepare(query)
	mensagemErro = "query_montagem_erro"
	CheckErro(w, r, mensagemErro, err)

	fmt.Println(query)

	_, err = stmt.Exec(novoPaciente.Nome, novoPaciente.Email, novoPaciente.Senha, novoPaciente.DataNascimento, novoPaciente.Hospital, novoPaciente.Carteira, "p", "a")
	mensagemErro = "query_exec_erro"
	CheckErro(w, r, mensagemErro, err)

	fmt.Println("teste")
	response, _ := json.Marshal(&novoPaciente)
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return
}

// Todos ...
func Todos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Chamando rota buscar todos pacientes ...")
	w.Header().Set("Content-Type", "application/json")

	pacientes = pacientes[:0]

	query := "SELECT codigo, UPPER(nome) , UPPER(email) ,data_nascimento, hospital ,carteira , ativo FROM paciente"
	rows, err := DB.Query(query)
	mensagemErro = "query_exec_erro"
	CheckErro(w, r, mensagemErro, err)

	for rows.Next() {
		paciente := paciente.Paciente{}
		rows.Scan(&paciente.Codigo, &paciente.Nome, &paciente.Email, &paciente.DataNascimento, &paciente.Hospital,
			&paciente.Carteira, &paciente.Ativo)
		pacientes = append(pacientes, paciente)
	}

	json.NewEncoder(w).Encode(pacientes)

}

// Buscar ...
func Buscar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Chamando rota buscar paciente ...")
	w.Header().Set("Content-Type", "application/json")

	pacientes = pacientes[:0]
	pacienteBuscar := paciente.Paciente{}

	err := json.NewDecoder(r.Body).Decode(&pacienteBuscar)
	mensagemErro = "erro_corpo"
	CheckErro(w, r, mensagemErro, err)

	query := "SELECT codigo, UPPER(nome) , UPPER(email) ,data_nascimento, hospital,carteira, ativo FROM paciente WHERE nome LIKE ?"
	rows, err := DB.Query(query, pacienteBuscar.Nome)
	mensagemErro = "query_exec_erro"
	CheckErro(w, r, mensagemErro, err)

	for rows.Next() {
		paciente := paciente.Paciente{}
		rows.Scan(&paciente.Codigo, &paciente.Nome, &paciente.Email, &paciente.DataNascimento, &paciente.Hospital,
			&paciente.Carteira, &paciente.Ativo)
		pacientes = append(pacientes, paciente)
	}

	json.NewEncoder(w).Encode(pacientes)

}

// Alterar ...
func Alterar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Chamando rota alterar paciente ...")
	w.Header().Set("Content-Type", "application/json")

	paciente := paciente.Paciente{}
	err := json.NewDecoder(r.Body).Decode(&paciente)
	mensagemErro = "erro_corpo"
	CheckErro(w, r, mensagemErro, err)

	// query := "UPDATE paciente SET nome = ? ,email = ?,data_nascimento = ?,hospital = ?,carteira = ? FROM paciente WHERE codigo = ?"
	stmt, err := DB.Prepare("UPDATE paciente SET nome = ? , email = ? , data_nascimento = ? , hospital = ? , carteira = ?, ativo = ? WHERE codigo = ?")
	mensagemErro = "query_exec_erro"
	CheckErro(w, r, mensagemErro, err)

	stmt.Exec(paciente.Nome, paciente.Email, paciente.DataNascimento, paciente.Hospital, paciente.Carteira, paciente.Ativo, paciente.Codigo)
	json.NewEncoder(w).Encode(paciente)

}

// CheckErro ...
func CheckErro(w http.ResponseWriter, r *http.Request, msg string, err error) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(map[string]interface{}{
			msg: err.Error(),
		})
		w.Write(response)
		return
	}
}
