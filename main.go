package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var _temp = template.Must(template.ParseGlob("./templates/*.html"))

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/{$}", homeHendler)
	router.HandleFunc("/", notFound)
	router.HandleFunc("/experience-list", ExperienceHendler)
	router.HandleFunc("/experience-details/{id}", ExperienceByHendler)

	fileServer := http.FileServer(http.Dir("./assets"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatalf("Erreur (Serveur) : Impossible de lancer le serveur : %v", err)
		return
	}
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	buffer := bytes.Buffer{}

	errRender := _temp.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		log.Printf("Erreur (Template) : Impossible de rendre le template %s : %v", name, errRender)
		http.Error(w, "Erreur chargement de la page", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buffer.WriteTo(w)
}

func homeHendler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "index", nil)
}

type Experience struct {
	Id         int
	Name       string
	LienGithub string
	Image      []string
}

var listExperience = []Experience{
	{1, "Projet Morpion", "https://github.com/MathieuJULIEN/Morpion-py", []string{"/static/image/morpion/morpion-en-jeu.png", "/static/image/morpion/rendu-global-code-morpion.png", "/static/image/morpion/start-morpion.png"}},
	{2, "Projet Crypt-Decrypt", "https://github.com/MathieuJULIEN/Micro-CryptDecrypt", []string{"/static/image/crypt-decrypt/1er-rendu-global-crypt-decrypt.png", "/static/image/crypt-decrypt/2eme-rendu-global-crypt-decrypt.png", "/static/image/crypt-decrypt/arboressance-crypt-decrypt.png", "/static/image/crypt-decrypt/fonctionnement-crypt.png", "/static/image/crypt-decrypt/fonctionnement-decrypt.png", "/static/image/crypt-decrypt/lancement-crypt-decrypt.png"}},
	{3, "Projet Jeu du Pendu", "https://github.com/kost84/TP-Jeu-du-pendu-CLI", []string{"/static/image/pendu/arboressance-pendu.png", "/static/image/pendu/fonctionnement-pendu.png", "/static/image/pendu/lancement-pendu.png", "/static/image/pendu/rendu-global-pendu.png"}},
}

func ExperienceHendler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "experience-list", listExperience)
}

func ExperienceByHendler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	id, errId := strconv.Atoi(r.PathValue("id"))
	if errId != nil {
		http.Error(w, "Erreur données invalides...", http.StatusBadRequest)
		return
	}

	for _, experience := range listExperience {
		if experience.Id == id {
			renderTemplate(w, "experience-details", experience)
			return
		}
	}
	http.Error(w, "Exeperience introuvable", http.StatusNotFound)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "error", nil)
}
